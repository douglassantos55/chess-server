package pkg

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
		Type: CreateGame,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "2s",
				Increment: "0s",
			},
		},
	})

	var params1 GameStart
	var params2 GameStart

	for params1.GameId == uuid.Nil || params2.GameId == uuid.Nil {
		select {
		case res := <-p1.Outgoing:
			params1 = res.Payload.(GameStart)
		case res := <-p2.Outgoing:
			params2 = res.Payload.(GameStart)
		case <-time.After(time.Second):
			t.Error("Expected game response, timedout instead")
		}
	}

	if params1.Color != White {
		t.Error("Expected p1 to be white")
	}
	if params2.Color != Black {
		t.Error("Expected p2 to be black")
	}

	if params1.TimeControl.Duration != "2s" {
		t.Errorf("Expected 2s duration, got %v", params1.TimeControl.Duration)
	}
	if params2.TimeControl.Duration != "2s" {
		t.Errorf("Expected 2s duration, got %v", params2.TimeControl.Duration)
	}

	game := gameManager.FindGame(params1.GameId)

	if game == nil {
		t.Error("Expected game to exist, got nil instead")
	}
}

func TestMovePieceHandler(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type: CreateGame,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "100ms",
				Increment: "0s",
			},
		},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameStart)
	game := gameManager.FindGame(params.GameId)

	go gameManager.Process(Message{
		Type: Move,
		Payload: map[string]interface{}{
			"from":    "e2",
			"to":      "e4",
			"game_id": game.Id.String(),
		},
	})

	<-p2.Outgoing

	if game.board.Square("e2") != Empty() {
		t.Errorf("Expected e2 to be empty, got %v", game.board.Square("e2"))
	}
	if game.board.Square("e4") != Pawn(White) {
		t.Errorf("Expected e4 to have a pawn, got %v", game.board.Square("e4"))
	}

	go gameManager.Process(Message{
		Type: Move,
		Payload: map[string]interface{}{
			"from":    "e7",
			"to":      "e5",
			"game_id": game.Id.String(),
		},
	})

	<-p1.Outgoing

	if game.board.Square("e7") != Empty() {
		t.Errorf("Expected e7 to be empty, got %v", game.board.Square("e7"))
	}
	if game.board.Square("e5") != Pawn(Black) {
		t.Errorf("Expected e5 to have a pawn, got %v", game.board.Square("e5"))
	}

	select {
	case <-time.After(time.Second):
		t.Error("Expected game over response, got timeout")
	case response := <-p1.Outgoing:
		result := response.Payload.(GameOverResponse)
		if result.Reason != "Timeout" {
			t.Errorf("Game should end with timeout")
		}
		if result.Winner {
			t.Error("Expected black to win on time")
		}
	case response := <-p2.Outgoing:
		result := response.Payload.(GameOverResponse)
		if result.Reason != "Timeout" {
			t.Errorf("Game should end with timeout")
		}
		if !result.Winner {
			t.Error("Expected black to win on time")
		}
	}
}

func TestSendsMoveEventToPlayer(t *testing.T) {
	manager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go manager.Process(Message{
		Type: CreateGame,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "4s",
				Increment: "0s",
			},
		},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameStart)

	go manager.Process(Message{
		Type: Move,
		Payload: map[string]interface{}{
			"from":    "e2",
			"to":      "e4",
			"game_id": params.GameId.String(),
		},
	})

	select {
	case <-time.After(time.Second):
		t.Error("Expected response, got timeout")
	case response := <-p2.Outgoing:
		if response.Type != StartTurn {
			t.Errorf("Expected StartTurn, got %v", response.Type)
		}

		payload := response.Payload.(MoveResponse)
		if payload.Time < time.Second {
			t.Errorf("Expected 1s, got %v", payload.Time)
		}
	}
}

func TestGameOver(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type: CreateGame,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "100ms",
				Increment: "0s",
			},
		},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	select {
	case res := <-p1.Outgoing:
		result := res.Payload.(GameOverResponse)
		if result.Winner {
			t.Error("Expected black to win on time")
		}
		if result.Reason != "Timeout" {
			t.Errorf("Expected black to win on time, got %v", result.Reason)
		}
	case res := <-p2.Outgoing:
		result := res.Payload.(GameOverResponse)
		if !result.Winner {
			t.Error("Expected black to win on time")
		}
		if result.Reason != "Timeout" {
			t.Errorf("Expected black to win on time, got %v", result.Reason)
		}
	}

	params := res.Payload.(GameStart)
	if gameManager.FindGame(params.GameId) != nil {
		t.Error("Expected game to be removed")
	}
}

func TestWhiteResign(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type: CreateGame,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "5s",
				Increment: "1s",
			},
		},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameStart)

	go gameManager.Process(Message{
		Player:  p1,
		Type:    Resign,
		Payload: params.GameId,
	})

	select {
	case res := <-p1.Outgoing:
		result := res.Payload.(GameOverResponse)
		if result.Winner {
			t.Error("Expected black to win by resignation")
		}
	case res := <-p2.Outgoing:
		result := res.Payload.(GameOverResponse)
		if !result.Winner {
			t.Error("Expected black to win by resignation")
		}
	case <-time.After(time.Second):
		t.Error("Expected game over, got timeout instead")
	}

	if gameManager.FindGame(params.GameId) != nil {
		t.Error("Expected game to be removed")
	}
}

func TestBlackResign(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type: CreateGame,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "3s",
				Increment: "0s",
			},
		},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameStart)

	go gameManager.Process(Message{
		Player:  p2,
		Type:    Resign,
		Payload: params.GameId,
	})

	select {
	case res := <-p1.Outgoing:
		result := res.Payload.(GameOverResponse)
		if !result.Winner {
			t.Error("Expected white to win by resignation")
		}
	case res := <-p2.Outgoing:
		result := res.Payload.(GameOverResponse)
		if result.Winner {
			t.Error("Expected white to win by resignation")
		}
	case <-time.After(time.Second):
		t.Error("Expected game over, got timeout instead")
	}

	if gameManager.FindGame(params.GameId) != nil {
		t.Error("Expected game to be removed")
	}
}

func TestBlackDisconnectEndsGame(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type: CreateGame,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "5s",
				Increment: "2s",
			},
		},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameStart)

	go gameManager.Process(Message{
		Player: p2,
		Type:   Disconnected,
	})

	select {
	case res := <-p1.Outgoing:
		result := res.Payload.(GameOverResponse)
		if !result.Winner {
			t.Error("Expected white to win by abandonment")
		}
		if result.Reason != "Abandonment" {
			t.Errorf("Expected white to win by abandonment, got %v", result.Reason)
		}
	case <-time.After(time.Second):
		t.Error("Expected game over, got timeout instead")
	}

	select {
	case <-time.After(time.Second):
	case <-p2.Outgoing:
		t.Error("Expected white to win by abandonment")
	}

	if gameManager.FindGame(params.GameId) != nil {
		t.Error("Expected game to be removed")
	}
}

func TestWhiteDisconnectEndsGame(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type: CreateGame,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "2s",
				Increment: "1s",
			},
		},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameStart)

	go gameManager.Process(Message{
		Player: p1,
		Type:   Disconnected,
	})

	select {
	case <-time.After(time.Second):
	case <-p1.Outgoing:
		t.Error("Disconnected player should not receive a response")
	}

	select {
	case res := <-p2.Outgoing:
		result := res.Payload.(GameOverResponse)
		if result.Reason != "Abandonment" {
			t.Errorf("Expected black to win by abandonment, got %v", result.Reason)
		}
		if !result.Winner {
			t.Error("Expected black to win by abandonment")
		}
	case <-time.After(time.Second):
		t.Error("Expected game over, got timeout instead")
	}

	if gameManager.FindGame(params.GameId) != nil {
		t.Error("Expected game to be removed")
	}
}
