package main

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func wait(cb func()) chan bool {
	channel := make(chan bool)

	go func() {
		cb()
		channel <- true
	}()

	return channel
}

func TestIgnoresIrrelevantEvents(t *testing.T) {
	gameManager := NewGameManager()

	<-wait(func() {
		gameManager.Process(Message{
			Type: QueueUp,
		})
	})

	if len(gameManager.games) != 0 {
		t.Errorf("Expected no games to exist, got %v instead", len(gameManager.games))
	}
}

func TestSendsResponseToPlayers(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type:    CreateGame,
		Payload: []*Player{p1, p2},
	})

	var params1 GameParams
	var params2 GameParams

	for params1.GameId == uuid.Nil || params2.GameId == uuid.Nil {
		select {
		case res1 := <-p1.Outgoing:
			params1 = res1.Payload.(GameParams)
		case res2 := <-p2.Outgoing:
			params2 = res2.Payload.(GameParams)
		case <-time.After(time.Second):
			t.Error("Expected game response, timedout instead")
		}
	}

	if params1.GameId != params2.GameId {
		t.Errorf("Expected same game ID, got %v and %v", params1.GameId, params2.GameId)
	}
}

func TestCreatesGame(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type:    CreateGame,
		Payload: []*Player{p1, p2},
	})

	var params1 GameParams
	var params2 GameParams

	for params1.GameId == uuid.Nil || params2.GameId == uuid.Nil {
		select {
		case res := <-p1.Outgoing:
			params1 = res.Payload.(GameParams)
		case res := <-p2.Outgoing:
			params2 = res.Payload.(GameParams)
		case <-time.After(time.Second):
			t.Error("Expected game response, timedout instead")
		}
	}

	game := gameManager.FindGame(params1.GameId)

	if game == nil {
		t.Error("Expected game to exist, got nil instead")
	}
}

func TestAssignsColorsToPlayers(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type:    CreateGame,
		Payload: []*Player{p1, p2},
	})

	var params1 GameParams
	var params2 GameParams

	for params1.GameId == uuid.Nil || params2.GameId == uuid.Nil {
		select {
		case res1 := <-p1.Outgoing:
			params1 = res1.Payload.(GameParams)
		case res2 := <-p2.Outgoing:
			params2 = res2.Payload.(GameParams)
		case <-time.After(time.Second):
			t.Error("Expected game response, timedout instead")
		}
	}

	if params1.Color == params2.Color {
		t.Errorf("Both players have the same color %v", params1.Color)
	}
}

func TestGameOverOnTime(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type:    CreateGame,
		Payload: []*Player{p1, p2},
	})

	<-p1.Outgoing
	<-p2.Outgoing

	time.Sleep(time.Second)

	select {
	case res := <-p1.Outgoing:
		if res.Type != GameOver {
			t.Errorf("Expected 'GameOver', got '%v'", res.Type)
		}
	case <-time.After(time.Second):
		t.Error("Timeout")
	}

	select {
	case res := <-p2.Outgoing:
		if res.Type != GameOver {
			t.Errorf("Expected 'GameOver', got '%v'", res.Type)
		}
	case <-time.After(time.Second):
		t.Error("Timeout")
	}
}
