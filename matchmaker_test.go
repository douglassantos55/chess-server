package main

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestIgnoresQueueUp(t *testing.T) {
	matchmaker := NewMatchMaker(time.Second)

	p1 := NewTestPlayer()

	go matchmaker.Process(Message{
		Type:    QueueUp,
		Payload: []*Player{p1},
	})

	select {
	case <-p1.Outgoing:
		t.Error("Should not handle QueueUp event")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestCreatesMatch(t *testing.T) {
	matchmaker := NewMatchMaker(time.Second)

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go matchmaker.Process(Message{
		Type:    MatchFound,
		Payload: []*Player{p1, p2},
	})

	res1 := <-p1.Outgoing
	res2 := <-p2.Outgoing

	matchId1 := res1.Payload.(uuid.UUID)
	matchId2 := res2.Payload.(uuid.UUID)

	if matchId1 != matchId2 {
		t.Errorf("Expected the same match ID, got %v and %v", matchId1, matchId2)
	}

	if !matchmaker.HasMatches() {
		t.Error("Expected matchmaker to have a match")
	}
}

func TestAsksForConfirmation(t *testing.T) {
	matchmaker := NewMatchMaker(time.Second)

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go matchmaker.Process(Message{
		Type:    MatchFound,
		Payload: []*Player{p1, p2},
	})

	select {
	case res := <-p1.Outgoing:
		if res.Type != ConfirmMatch {
			t.Errorf("Expected confirm match, got %v", res.Type)
		}
	case <-time.After(time.Second):
		t.Error("Expected response, got timeout instead")
	}

	select {
	case res := <-p2.Outgoing:
		if res.Type != ConfirmMatch {
			t.Errorf("Expected confirm match, got %v", res.Type)
		}
	case <-time.After(time.Second):
		t.Error("Expected response, got timeout instead")
	}
}

func TestConcurrency(t *testing.T) {
	matchmaker := NewMatchMaker(time.Second)

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()
	p3 := NewTestPlayer()

	go matchmaker.Process(Message{
		Type:    MatchFound,
		Payload: []*Player{p1},
	})
	go matchmaker.Process(Message{
		Type:    MatchFound,
		Payload: []*Player{p2},
	})
	go matchmaker.Process(Message{
		Type:    MatchFound,
		Payload: []*Player{p3},
	})

	responses := []Response{}

	for len(responses) != 3 {
		select {
		case res := <-p1.Outgoing:
			responses = append(responses, res)
		case res := <-p2.Outgoing:
			responses = append(responses, res)
		case res := <-p3.Outgoing:
			responses = append(responses, res)
		case <-time.After(time.Second):
			t.Error("Expected responses from server, got timeout instead")
		}
	}
}

func TestRequeuesConfirmedAfterTimeout(t *testing.T) {
	defer func() {
		Dispatcher = nil
	}()

	Dispatcher = make(chan Message)
	matchmaker := NewMatchMaker(200 * time.Millisecond)

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go matchmaker.Process(Message{
		Type:    MatchFound,
		Payload: []*Player{p1, p2},
	})

	var response Response

	select {
	case res := <-p1.Outgoing:
		response = res
	case res := <-p2.Outgoing:
		response = res
	case <-time.After(time.Second):
		t.Error("Expected response, got timeout")
	}

	matchId := response.Payload.(uuid.UUID)

	go matchmaker.Process(Message{
		Player:  p1,
		Payload: matchId,
		Type:    MatchConfirmed,
	})

	res := <-p1.Outgoing
	if res.Type != WaitOtherPlayers {
		t.Errorf("Expected wait other players, got %v", res.Type)
	}

	time.Sleep(200 * time.Millisecond)

	<-p1.Outgoing

	select {
	case queueUp := <-Dispatcher:
		if queueUp.Type != QueueUp {
			t.Error("Expected confirmed to be requeued", queueUp.Type)
		}
	case <-time.After(time.Second):
		t.Error("Expected response, got timeout")
	}
}

func TestCancelsMatchIfNoConfirmation(t *testing.T) {
	matchmaker := NewMatchMaker(200 * time.Millisecond)
	p1 := NewTestPlayer()

	go matchmaker.Process(Message{
		Type:    MatchFound,
		Payload: []*Player{p1},
	})

	<-p1.Outgoing
	time.Sleep(300 * time.Millisecond)

	if matchmaker.HasMatches() {
		t.Error("Expected match to be canceled")
	}
}

func TestDispatchesGameStart(t *testing.T) {
	defer func() {
		Dispatcher = nil
	}()

	Dispatcher = make(chan Message)
	matchmaker := NewMatchMaker(time.Second)

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go matchmaker.Process(Message{
		Type:    MatchFound,
		Payload: []*Player{p1, p2},
	})

	response := <-p1.Outgoing
	<-p2.Outgoing

	matchId := response.Payload.(uuid.UUID)

	go matchmaker.Process(Message{
		Player:  p1,
		Payload: matchId,
		Type:    MatchConfirmed,
	})
	go matchmaker.Process(Message{
		Player:  p2,
		Payload: matchId,
		Type:    MatchConfirmed,
	})

	waitP1 := <-p1.Outgoing
	if waitP1.Type != WaitOtherPlayers {
		t.Errorf("Expected wait other players, got %v", waitP1.Type)
	}

	waitP2 := <-p2.Outgoing
	if waitP2.Type != WaitOtherPlayers {
		t.Errorf("Expected wait other players, got %v", waitP2.Type)
	}

	select {
	case res := <-Dispatcher:
		if res.Type != GameStart {
			t.Errorf("Expected game start, got %v", res.Type)
		}

		players := res.Payload.([]*Player)

		if len(players) != MAX_PLAYERS {
			t.Errorf("Expected 2 players, got %v", len(players))
		}
	case <-time.After(time.Second):
		t.Error("Expected response, got timeout")
	}
}

func TestRefuseMatch(t *testing.T) {
	defer func() {
		Dispatcher = nil
	}()

	Dispatcher = make(chan Message)
	matchmaker := NewMatchMaker(time.Second)

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go matchmaker.Process(Message{
		Type:    MatchFound,
		Payload: []*Player{p1, p2},
	})

	response := <-p1.Outgoing
	<-p2.Outgoing

	matchId := response.Payload.(uuid.UUID)

	go matchmaker.Process(Message{
		Player:  p1,
		Payload: matchId,
		Type:    MatchConfirmed,
	})

	waitP1 := <-p1.Outgoing
	if waitP1.Type != WaitOtherPlayers {
		t.Errorf("Expected wait other players, got %v", waitP1.Type)
	}

	go matchmaker.Process(Message{
		Player:  p2,
		Payload: matchId,
		Type:    MatchDeclined,
	})

	notification := <-p1.Outgoing
	if notification.Type != MatchCanceled {
		t.Errorf("Expected match canceled response, got %v", notification.Type)
	}

	select {
	case queueUp := <-Dispatcher:
		if queueUp.Type != QueueUp {
			t.Error("Expected confirmed to be requeued", queueUp.Type)
		}
	case <-time.After(time.Second):
		t.Error("Expected response, got timeout")
	}

	if matchmaker.HasMatches() {
		t.Error("Expected match to be canceled")
	}
}
