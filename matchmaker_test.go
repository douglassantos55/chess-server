package main

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestIgnoresQueueUp(t *testing.T) {
	matchmaker := NewMatchMaker()

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
	matchmaker := NewMatchMaker()

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
	matchmaker := NewMatchMaker()

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
	matchmaker := NewMatchMaker()

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
