package main

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAssignsColorsToPlayers(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame(500*time.Millisecond, []*Player{p1, p2})
	go game.Start()

	var params1 GameParams
	var params2 GameParams

	for params1.GameId == uuid.Nil || params2.GameId == uuid.Nil {
		select {
		case res1 := <-p1.Outgoing:
			params1 = res1.Payload.(GameParams)
		case res2 := <-p2.Outgoing:
			params2 = res2.Payload.(GameParams)
		case <-time.After(500 * time.Millisecond):
			t.Error("Expected game response, timedout instead")
		}
	}

	if params1.Color == params2.Color {
		t.Errorf("Both players have the same color %v", params1.Color)
	}

	<-game.Over
}

func TestPausesTimerOnEndTurn(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame(500*time.Millisecond, []*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing // StartGame
	<-p2.Outgoing // StartGame

	game.EndTurn()

	select {
	case <-time.After(500 * time.Millisecond):
	case <-game.Over:
		t.Error("Should not end game on time")
	}
}

func TestStartsWhiteTimer(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame(500*time.Millisecond, []*Player{p1, p2})
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
	case <-time.After(600 * time.Millisecond):
		t.Error("Timeout")
	}
}

func TestTimerStops(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame(500*time.Millisecond, []*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing
	<-p2.Outgoing

	time.Sleep(200 * time.Millisecond)
	game.EndTurn()

	if game.Current.Next.left < 200*time.Millisecond {
		t.Errorf("Time left should be 250ms, got %v", game.Current.Next.left)
	}
}

func TestTimerContinues(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame(500*time.Millisecond, []*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing
	<-p2.Outgoing

	time.Sleep(250 * time.Millisecond)

	game.EndTurn() // stops white and moves current to black
	game.EndTurn() // stops black and moves current back to white

	game.StartTurn()

	select {
	case <-game.Over:
	case <-time.After(500 * time.Millisecond):
		t.Error("Expected game over within 250ms, got timeout after 500ms")
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
	game.EndTurn()
	game.StartTurn()

	select {
	case result := <-game.Over:
		if result.Winner != p1 {
			t.Error("Expected white to win on time")
		}
	case <-time.After(600 * time.Millisecond):
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

	game.EndTurn() // end white's turn

	game.StartTurn() // start black's turn
	game.Move("e7", "e5")
	game.EndTurn()

	game.StartTurn() // start white's turn

	select {
	case result := <-game.Over:
		if result.Winner != p2 {
			t.Error("Expected black to win on time")
		}
	case <-time.After(600 * time.Millisecond):
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

	if game.Move("e7", "e5") {
		t.Error("White should not be able to move black's pieces")
	}

	select {
	case result := <-game.Over:
		if result.Winner != p2 {
			t.Error("Expected black to win on time")
		}
	case <-time.After(600 * time.Millisecond):
		t.Error("Expected white's timer to run out")
	}
}

func TestBlackCannotMoveWhitePiece(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame(500*time.Millisecond, []*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing
	<-p2.Outgoing

	if game.Move("e2", "e4") { // white's turn
		game.EndTurn()
		game.StartTurn()
	}

	if game.Move("d2", "d4") { // black's turn
		t.Error("Black should not be able to move white's pieces")
	}

	select {
	case result := <-game.Over:
		if result.Winner != p1 {
			t.Error("Expected white to win on time")
		}
	case <-time.After(600 * time.Millisecond):
		t.Error("Expected black's timer to run out")
	}
}

func TestCheckmate(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame(5*time.Second, []*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing
	<-p2.Outgoing

	game.Move("f2", "f3") // white's turn
	game.EndTurn()

	game.Move("e7", "e5") // black's turn
	game.EndTurn()

	game.Move("g2", "g4") // white's turn
	game.EndTurn()

	game.Move("d8", "h4") // black's turn
	game.EndTurn()

	if !game.IsCheckmate() {
		t.Error("Expected black to win by checkmate")
	}

	go game.Checkmate()

	select {
	case result := <-game.Over:
		if result.Winner != p2 {
			t.Error("Expected black to win")
		}
		if result.Reason != "Checkmate" {
			t.Errorf("Expected black to win by checkmate, got %v", result.Reason)
		}
	case <-time.After(5 * time.Second):
		t.Error("Expected black to win by checkmate, timeout instead")
	}
}

func TestBlockCheckmate(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	game := NewGame(time.Second, []*Player{p1, p2})
	go game.Start()

	<-p1.Outgoing
	<-p2.Outgoing

	game.Move("e2", "e4") // white's turn
	game.EndTurn()

	game.StartTurn()
	game.Move("e7", "e6") // black's turn
	game.EndTurn()

	game.StartTurn()
	game.Move("b2", "b3") // white's turn
	game.EndTurn()

	game.StartTurn()
	game.Move("d8", "h4") // black's turn
	game.EndTurn()

	game.StartTurn()
	game.Move("h2", "h3") // white's turn
	game.EndTurn()

	game.StartTurn()
	game.Move("h4", "e4") // black's turn
	game.EndTurn()

	game.StartTurn() // white's turn

	if game.IsCheckmate() {
		t.Error("Should not be checkmate, can block on e2")
	}

	select {
	case result := <-game.Over:
		if result.Reason == "Checkmate" {
			t.Error("Game should not end with checkmate, queen/bishop can block on e2")
		}
		if result.Winner != p2 {
			t.Error("Expected black to win on time")
		}
	}
}
