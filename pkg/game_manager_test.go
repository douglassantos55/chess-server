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

func TestMovePieceHandler(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type:    CreateGame,
		Payload: []*Player{p1, p2},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameParams)
	game := gameManager.FindGame(params.GameId)

	go gameManager.Process(Message{
		Type: Move,
		Payload: MovePiece{
			From:   "e2",
			To:     "e4",
			GameId: game.Id,
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
		Payload: MovePiece{
			From:   "e7",
			To:     "e5",
			GameId: game.Id,
		},
	})

	<-p1.Outgoing

	if game.board.Square("e7") != Empty() {
		t.Errorf("Expected e7 to be empty, got %v", game.board.Square("e7"))
	}
	if game.board.Square("e5") != Pawn(Black) {
		t.Errorf("Expected e5 to have a pawn, got %v", game.board.Square("e5"))
	}

	assertWinner := func(result GameResult) {
		if result.Reason != "Timeout" {
			t.Errorf("Game should end with timeout")
		}
		if result.Winner != p2 {
			t.Error("Expected black to win on time")
		}
	}

	select {
	case <-time.After(2 * time.Second):
		t.Error("Expected game over response, got timeout")
	case response := <-p1.Outgoing:
		result := response.Payload.(GameResult)
		assertWinner(result)
	case response := <-p2.Outgoing:
		result := response.Payload.(GameResult)
		assertWinner(result)
	}
}

func TestSendsMoveEventToPlayer(t *testing.T) {
	manager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go manager.Process(Message{
		Type:    CreateGame,
		Payload: []*Player{p1, p2},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameParams)

	go manager.Process(Message{
		Type: Move,
		Payload: MovePiece{
			From:   "e2",
			To:     "e4",
			GameId: params.GameId,
		},
	})

	select {
	case <-time.After(time.Second):
		t.Error("Expected response, got timeout")
	case response := <-p2.Outgoing:
		if response.Type != StartTurn {
			t.Errorf("Expected StartTurn, got %v", response.Type)
		}

		if timer := response.Payload.(time.Duration); timer < time.Second {
			t.Errorf("Expected 1s, got %v", timer)
		}
	}
}

func TestGameOver(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type:    CreateGame,
		Payload: []*Player{p1, p2},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameParams)

	select {
	case res := <-p1.Outgoing:
		result := res.Payload.(GameResult)
		if result.Winner != p2 {
			t.Error("Expected black to win on time")
		}
	case res := <-p2.Outgoing:
		result := res.Payload.(GameResult)
		if result.Winner != p2 {
			t.Error("Expected black to win on time")
		}
	}

	if gameManager.FindGame(params.GameId) != nil {
		t.Error("Expected game to be removed")
	}
}

func TestWhiteResign(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type:    CreateGame,
		Payload: []*Player{p1, p2},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameParams)

	go gameManager.Process(Message{
		Player:  p1,
		Type:    Resign,
		Payload: params.GameId,
	})

	select {
	case res := <-p1.Outgoing:
		result := res.Payload.(GameResult)
		if result.Winner != p2 {
			t.Error("Expected black to win by resignation")
		}
	case res := <-p2.Outgoing:
		result := res.Payload.(GameResult)
		if result.Winner != p2 {
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
		Type:    CreateGame,
		Payload: []*Player{p1, p2},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameParams)

	go gameManager.Process(Message{
		Player:  p2,
		Type:    Resign,
		Payload: params.GameId,
	})

	select {
	case res := <-p1.Outgoing:
		result := res.Payload.(GameResult)
		if result.Winner != p1 {
			t.Error("Expected white to win by resignation")
		}
	case res := <-p2.Outgoing:
		result := res.Payload.(GameResult)
		if result.Winner != p1 {
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
		Type:    CreateGame,
		Payload: []*Player{p1, p2},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameParams)

	go gameManager.Process(Message{
		Player: p2,
		Type:   Disconnected,
	})

	select {
	case res := <-p1.Outgoing:
		result := res.Payload.(GameResult)
		if result.Winner != p1 {
			t.Error("Expected white to win by resignation")
		}
	case res := <-p2.Outgoing:
		result := res.Payload.(GameResult)
		if result.Winner != p1 {
			t.Error("Expected white to win by resignation")
		}
	case <-time.After(time.Second):
		t.Error("Expected game over, got timeout instead")
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
		Type:    CreateGame,
		Payload: []*Player{p1, p2},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameParams)

	go gameManager.Process(Message{
		Player: p1,
		Type:   Disconnected,
	})

	select {
	case res := <-p1.Outgoing:
		result := res.Payload.(GameResult)
		if result.Winner != p2 {
			t.Error("Expected black to win by resignation")
		}
	case res := <-p2.Outgoing:
		result := res.Payload.(GameResult)
		if result.Winner != p2 {
			t.Error("Expected black to win by resignation")
		}
	case <-time.After(time.Second):
		t.Error("Expected game over, got timeout instead")
	}

	if gameManager.FindGame(params.GameId) != nil {
		t.Error("Expected game to be removed")
	}
}
