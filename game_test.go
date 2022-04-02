package main

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAssignsColorsToPlayers(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame(time.Second, []*Player{p1, p2})
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

	game := NewGame(time.Second, []*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing // StartGame
	<-p2.Outgoing // StartGame

	game.EndTurn()

	select {
	case <-time.After(time.Second):
	case <-game.Over:
		t.Error("Should not end game on time")
	}
}

func TestStartsWhiteTimer(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame(time.Second, []*Player{p1, p2})
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

	game := NewGame(time.Second, []*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing
	<-p2.Outgoing

	time.Sleep(500 * time.Millisecond)
	game.EndTurn()

	if game.Current.Next.left >= time.Second {
		t.Errorf("Time left should be 500ms, got %v", game.Current.Next.left)
	}
}

func TestTimerContinues(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame(time.Second, []*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing
	<-p2.Outgoing

	time.Sleep(500 * time.Millisecond)

	game.EndTurn() // stops white and moves current to black
	game.EndTurn() // stops black and moves current back to white

	game.StartTurn()

	select {
	case <-game.Over:
	case <-time.After(time.Second):
		t.Error("Expected game over within 500ms, got timeout after 1s")
	}
}

func TestWhiteMovePassesTurn(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame(500*time.Millisecond, []*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing
	<-p2.Outgoing

	game.Move("e2", "e4")

	select {
	case result := <-game.Over:
		if result.Winner != p1 {
			t.Error("Expected white to win on time")
		}
	case <-time.After(time.Second):
		t.Error("Expected game over")
	}
}

func TestBlackMovePassesTurn(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame(500*time.Millisecond, []*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing
	<-p2.Outgoing

	game.EndTurn()
	game.StartTurn()

	game.Move("e7", "e5")

	select {
	case result := <-game.Over:
		if result.Winner != p2 {
			t.Error("Expected black to win on time")
		}
	case <-time.After(time.Second):
		t.Error("Expected game over")
	}
}

func TestWhiteCannotMoveBlackPiece(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame(500*time.Millisecond, []*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing
	<-p2.Outgoing

	game.Move("e7", "e5")

	select {
	case result := <-game.Over:
		if result.Winner != p2 {
			t.Error("Expected black to win on time")
		}
	case <-time.After(time.Second):
		t.Error("Expected white's timer to run out")
	}
}
