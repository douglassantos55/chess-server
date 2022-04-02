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
