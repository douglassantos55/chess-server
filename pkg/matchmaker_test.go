package pkg

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
		Type: MatchFound,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "10m",
				Increment: "0s",
			},
		},
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
		Type: MatchFound,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "1m",
				Increment: "0s",
			},
		},
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
		Type: MatchFound,
		Payload: MatchParams{
			Players: []*Player{p1},
			TimeControl: TimeControl{
				Duration:  "1m",
				Increment: "0s",
			},
		},
	})
	go matchmaker.Process(Message{
		Type: MatchFound,
		Payload: MatchParams{
			Players: []*Player{p2},
			TimeControl: TimeControl{
				Duration:  "1m",
				Increment: "0s",
			},
		},
	})
	go matchmaker.Process(Message{
		Type: MatchFound,
		Payload: MatchParams{
			Players: []*Player{p3},
			TimeControl: TimeControl{
				Duration:  "1m",
				Increment: "0s",
			},
		},
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
	matchmaker := NewMatchMaker(200 * time.Millisecond)

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go matchmaker.Process(Message{
		Type: MatchFound,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "10m",
				Increment: "0s",
			},
		},
	})

	var response1 Response
	var response2 Response

	for response1.Type == "" || response2.Type == "" {
		select {
		case res := <-p1.Outgoing:
			response1 = res
		case res := <-p2.Outgoing:
			response2 = res
		case <-time.After(time.Second):
			t.Error("Expected response, got timeout")
		}
	}

	matchId := response1.Payload.(uuid.UUID)

	go matchmaker.Process(Message{
		Player:  p1,
		Payload: matchId.String(),
		Type:    MatchConfirmed,
	})

	res := <-p1.Outgoing
	if res.Type != WaitOtherPlayers {
		t.Errorf("Expected wait other players, got %v", res.Type)
	}

	time.Sleep(200 * time.Millisecond)

	// match canceled response
	<-p1.Outgoing
	<-p2.Outgoing

	select {
	case queueUp := <-Dispatcher:
		if queueUp.Type != QueueUp {
			t.Error("Expected confirmed to be requeued", queueUp.Type)
		}

		payload := queueUp.Payload.(map[string]interface{})
		if payload["duration"] != "10m" {
			t.Errorf("Expected 10m duration, got %v", payload["duration"])
		}
		if payload["increment"] != "0s" {
			t.Errorf("Expected 0s increment, got %v", payload["increment"])
		}
	case <-time.After(time.Second):
		t.Error("Expected response, got timeout")
	}
}

func TestCancelsMatchIfNoConfirmation(t *testing.T) {
	matchmaker := NewMatchMaker(200 * time.Millisecond)
	p1 := NewTestPlayer()

	go matchmaker.Process(Message{
		Type: MatchFound,
		Payload: MatchParams{
			Players: []*Player{p1},
			TimeControl: TimeControl{
				Duration:  "5m",
				Increment: "5s",
			},
		},
	})

	<-p1.Outgoing
	time.Sleep(300 * time.Millisecond)

	if matchmaker.HasMatches() {
		t.Error("Expected match to be canceled")
	}
}

func TestDispatchesGameStart(t *testing.T) {
	matchmaker := NewMatchMaker(time.Second)

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go matchmaker.Process(Message{
		Type: MatchFound,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "5m",
				Increment: "0s",
			},
		},
	})

	response := <-p1.Outgoing
	<-p2.Outgoing

	matchId := response.Payload.(uuid.UUID)

	go matchmaker.Process(Message{
		Player:  p1,
		Payload: matchId.String(),
		Type:    MatchConfirmed,
	})

	waitP1 := <-p1.Outgoing
	if waitP1.Type != WaitOtherPlayers {
		t.Errorf("Expected wait other players, got %v", waitP1.Type)
	}

	go matchmaker.Process(Message{
		Player:  p2,
		Payload: matchId.String(),
		Type:    MatchConfirmed,
	})

	waitP2 := <-p2.Outgoing
	if waitP2.Type != WaitOtherPlayers {
		t.Errorf("Expected wait other players, got %v", waitP2.Type)
	}

	select {
	case res := <-Dispatcher:
		if res.Type != CreateGame {
			t.Errorf("Expected game start, got %v", res.Type)
		}

		params := res.Payload.(MatchParams)

		if len(params.Players) != MAX_PLAYERS {
			t.Errorf("Expected 2 players, got %v", len(params.Players))
		}
		if params.TimeControl.Duration != "5m" {
			t.Errorf("Expected 5m duration, got %v", params.TimeControl.Duration)
		}
		if params.TimeControl.Increment != "0s" {
			t.Errorf("Expected 0s increment, got %v", params.TimeControl.Increment)
		}
	case <-time.After(time.Second):
		t.Error("Expected response, got timeout")
	}

	if matchmaker.HasMatches() {
		t.Error("Expected match to be canceled")
	}
}

func TestRefuseMatch(t *testing.T) {
	matchmaker := NewMatchMaker(time.Second)

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go matchmaker.Process(Message{
		Type: MatchFound,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "15m",
				Increment: "5s",
			},
		},
	})

	var response1 Response
	var response2 Response

	for response1.Type == "" || response2.Type == "" {
		select {
		case response := <-p1.Outgoing:
			response1 = response
		case response := <-p2.Outgoing:
			response2 = response
		}
	}

	matchId := response1.Payload.(uuid.UUID)

	go matchmaker.Process(Message{
		Player:  p1,
		Payload: matchId.String(),
		Type:    MatchConfirmed,
	})

	waitP1 := <-p1.Outgoing
	if waitP1.Type != WaitOtherPlayers {
		t.Errorf("Expected wait other players, got %v", waitP1.Type)
	}

	go matchmaker.Process(Message{
		Player:  p2,
		Payload: matchId.String(),
		Type:    MatchDeclined,
	})

	notification := <-p1.Outgoing
	if notification.Type != MatchCanceled {
		t.Errorf("Expected match canceled response, got %v", notification.Type)
	}

	<-p2.Outgoing

	select {
	case queueUp := <-Dispatcher:
		if queueUp.Type != QueueUp {
			t.Error("Expected confirmed to be requeued", queueUp.Type)
		}

		payload := queueUp.Payload.(map[string]interface{})
		if payload["duration"] != "15m" {
			t.Errorf("Expected 15m duration, got %v", payload["duration"])
		}
		if payload["increment"] != "5s" {
			t.Errorf("Expected 5s increment, got %v", payload["increment"])
		}
	case <-time.After(time.Second):
		t.Error("Expected response, got timeout")
	}

	if matchmaker.HasMatches() {
		t.Error("Expected match to be canceled")
	}
}

func TestDisconnect(t *testing.T) {
	matchmaker := NewMatchMaker(time.Second)

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go matchmaker.Process(Message{
		Type: MatchFound,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "1m",
				Increment: "1s",
			},
		},
	})

	var response1 Response
	var response2 Response

	for response1.Type == "" || response2.Type == "" {
		select {
		case response := <-p1.Outgoing:
			response1 = response
		case response := <-p2.Outgoing:
			response2 = response
		}
	}

	matchId := response1.Payload.(uuid.UUID)

	go matchmaker.Process(Message{
		Player:  p1,
		Payload: matchId.String(),
		Type:    MatchConfirmed,
	})

	waitP1 := <-p1.Outgoing
	if waitP1.Type != WaitOtherPlayers {
		t.Errorf("Expected wait other players, got %v", waitP1.Type)
	}

	go matchmaker.Process(Message{
		Type:   Disconnected,
		Player: p2,
	})

	canceled := <-p1.Outgoing
	if canceled.Type != MatchCanceled {
		t.Errorf("Expected match canceled, got %v", canceled.Type)
	}

	if matchmaker.HasMatches() {
		t.Error("Expected match to be canceled")
	}

	<-p2.Outgoing

	select {
	case queueUp := <-Dispatcher:
		if queueUp.Type != QueueUp {
			t.Error("Expected confirmed to be requeued", queueUp.Type)
		}

		payload := queueUp.Payload.(map[string]interface{})
		if payload["duration"] != "1m" {
			t.Errorf("Expected 1m duration, got %v", payload["duration"])
		}
		if payload["increment"] != "1s" {
			t.Errorf("Expected 1s increment, got %v", payload["increment"])
		}
	case <-time.After(time.Second):
		t.Error("Expected response, got timeout")
	}

}
