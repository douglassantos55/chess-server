package main

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAssignsColorsToPlayers(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame([]*Player{p1, p2})
	go game.Start()

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

func TestPausesTimerOnEndTurn(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame([]*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing // StartGame
	<-p2.Outgoing // StartGame

	game.EndTurn(White)

	select {
	case <-time.After(time.Second):
	case <-game.Over:
		t.Error("Should not end game on time")
	}
}

func TestGameOverOnTime(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame([]*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing
	<-p2.Outgoing

	select {
	case gameOver := <-game.Over:
		if gameOver.Winner != p2 {
			t.Errorf("Expected p2 to win, got %v", gameOver.Winner)
		}
		if gameOver.Loser != p1 {
			t.Errorf("Expected p1 to win, got %v", gameOver.Loser)
		}
	case <-time.After(2 * time.Second):
		t.Error("Timeout")
	}
}

func TestTimerStops(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame([]*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing
	<-p2.Outgoing

	time.Sleep(500 * time.Millisecond)
	game.EndTurn(White)

	if game.White.left >= time.Second {
		t.Errorf("Time left should be 500ms, got %v", game.White.left)
	}
}

func TestTimerContinues(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame([]*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing
	<-p2.Outgoing

	time.Sleep(500 * time.Millisecond)

	game.EndTurn(White)
	game.StartTurn(White)

	select {
	case <-game.Over:
	case <-time.After(time.Second):
		t.Error("Expected game over within 500ms, got timeout after 1s")
	}
}
