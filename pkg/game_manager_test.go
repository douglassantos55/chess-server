package pkg

import (
	"reflect"
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
				Duration:  "1s",
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
		if payload.To != "e4" {
			t.Errorf("Expected to e4, got %v", payload.To)
		}
		if payload.From != "e2" {
			t.Errorf("Expected from e2, got %v", payload.From)
		}
		if payload.GameId != params.GameId {
			t.Errorf("Expected game ID to be %v, got %v", params.GameId, payload.GameId)
		}
		if payload.Time < time.Second.Milliseconds() {
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
		Payload: params.GameId.String(),
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
		Payload: params.GameId.String(),
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

func TestUnknownProblem(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type: CreateGame,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "5m",
				Increment: "1s",
			},
		},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameStart)
	game := gameManager.FindGame(params.GameId)

	game.board.matrix[6]['f'] = Empty()
	game.board.matrix[6]['g'] = Empty()
	game.board.matrix[6]['d'] = Empty()
	game.board.matrix[6]['e'] = Rook(White)
	game.board.matrix[5]['d'] = Pawn(Black)
	game.board.matrix[5]['e'] = Rook(White)
	game.board.matrix[4]['d'] = Pawn(White)
	game.board.matrix[3]['c'] = Bishop(White)
	game.board.matrix[2]['c'] = Queen(Black)
	game.board.matrix[7]['b'] = Empty()
	game.board.matrix[7]['c'] = Empty()
	game.board.matrix[7]['d'] = Empty()
	game.board.matrix[7]['e'] = Empty()
	game.board.matrix[7]['f'] = Empty()
	game.board.matrix[7]['g'] = Empty()
	game.board.matrix[4]['f'] = Bishop(Black)
	game.board.matrix[4]['g'] = King(Black)
	game.board.matrix[0]['a'] = Empty()
	game.board.matrix[0]['b'] = Empty()
	game.board.matrix[0]['c'] = Empty()
	game.board.matrix[0]['d'] = Empty()
	game.board.matrix[0]['e'] = Empty()
	game.board.matrix[0]['f'] = Empty()
	game.board.matrix[0]['h'] = Empty()
	game.board.matrix[0]['g'] = King(White)
	game.board.matrix[1]['b'] = Empty()
	game.board.matrix[1]['c'] = Empty()
	game.board.matrix[1]['d'] = Empty()
	game.board.matrix[1]['e'] = Empty()

	game.Current.King = "g1"
	game.Current.Next.King = "g5"

	go gameManager.Process(Message{
		Type: Move,
		Payload: MovePiece{
			To:     "g7",
			From:   "e7",
			GameId: game.Id.String(),
		},
	})

	<-p2.Outgoing
	if game.board.Square("g7") != Rook(White) {
		t.Errorf("Expected rook on g7, got %v", game.board.Square("g7"))
	}
}

func TestYetAnotherUnknownProblem(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type: CreateGame,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "5m",
				Increment: "1s",
			},
		},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameStart)
	game := gameManager.FindGame(params.GameId)

	game.board.matrix[0]['a'] = Empty()
	game.board.matrix[0]['b'] = Empty()
	game.board.matrix[0]['c'] = King(White)
	game.board.matrix[0]['d'] = Knight(Black)
	game.board.matrix[0]['e'] = Empty()
	game.board.matrix[0]['f'] = Empty()
	game.board.matrix[0]['g'] = Empty()
	game.board.matrix[0]['h'] = Empty()

	game.board.matrix[1]['c'] = Empty()
	game.board.matrix[2]['c'] = Pawn(White)
	game.board.matrix[1]['d'] = Empty()
	game.board.matrix[3]['d'] = Pawn(White)
	game.board.matrix[1]['e'] = Empty()
	game.board.matrix[2]['e'] = Pawn(White)

	game.board.matrix[1]['f'] = Empty()
	game.board.matrix[5]['b'] = Bishop(White)
	game.board.matrix[2]['f'] = Knight(White)
	game.board.matrix[3]['f'] = Bishop(White)

	game.board.matrix[7]['a'] = Empty()
	game.board.matrix[7]['b'] = Empty()
	game.board.matrix[7]['c'] = Empty()
	game.board.matrix[7]['d'] = Empty()
	game.board.matrix[7]['e'] = Empty()
	game.board.matrix[7]['f'] = Empty()
	game.board.matrix[7]['g'] = Empty()
	game.board.matrix[7]['h'] = Queen(White)

	game.board.matrix[6]['c'] = Empty()
	game.board.matrix[6]['d'] = King(Black)
	game.board.matrix[6]['e'] = Empty()
	game.board.matrix[4]['a'] = Bishop(Black)
	game.board.matrix[4]['d'] = Pawn(Black)
	game.board.matrix[5]['e'] = Bishop(Black)

	game.board.matrix[6]['h'] = Empty()
	game.board.matrix[5]['h'] = Pawn(Black)

	game.Current.King = "c1"
	game.Current.Next.King = "d7"

	if len(game.board.IsThreatened("d7", White)) == 0 {
		t.Error("Black should be checked")
	}

	game.EndTurn()

	go gameManager.Process(Message{
		Type: Move,
		Payload: MovePiece{
			To:     "c6",
			From:   "d7",
			GameId: game.Id.String(),
		},
	})

	<-p1.Outgoing
	if !reflect.DeepEqual(game.board.Square("c6"), King(Black)) {
		t.Errorf("Expected king on c6, got %v", game.board.Square("c6"))
	}
}

func TestOneMoreUnknownProblem(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type: CreateGame,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "5m",
				Increment: "1s",
			},
		},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameStart)
	game := gameManager.FindGame(params.GameId)

	game.board.matrix[0]['a'] = Empty()
	game.board.matrix[0]['b'] = Empty()
	game.board.matrix[0]['c'] = Empty()
	game.board.matrix[0]['d'] = Empty()
	game.board.matrix[0]['e'] = Empty()
	game.board.matrix[0]['f'] = Empty()
	game.board.matrix[0]['g'] = King(White)
	game.board.matrix[0]['h'] = Empty()

	game.board.matrix[1]['a'] = Empty()
	game.board.matrix[1]['b'] = Empty()
	game.board.matrix[1]['c'] = Empty()
	game.board.matrix[1]['d'] = Empty()
	game.board.matrix[1]['e'] = Bishop(White)
	game.board.matrix[1]['f'] = Empty()
	game.board.matrix[1]['g'] = Empty()

	game.board.matrix[2]['a'] = Pawn(White)

	game.board.matrix[3]['c'] = Pawn(Black)
	game.board.matrix[3]['e'] = Pawn(White)
	game.board.matrix[3]['g'] = Pawn(White)

	game.board.matrix[4]['f'] = Pawn(White)
	game.board.matrix[4]['g'] = Pawn(Black)

	game.board.matrix[5]['a'] = King(Black)
	game.board.matrix[5]['b'] = Pawn(Black)
	game.board.matrix[5]['c'] = Knight(Black)
	game.board.matrix[5]['h'] = Pawn(Black)

	game.board.matrix[6]['a'] = Empty()
	game.board.matrix[6]['b'] = Empty()
	game.board.matrix[6]['c'] = Pawn(Black)
	game.board.matrix[6]['d'] = Empty()
	game.board.matrix[6]['e'] = Empty()
	game.board.matrix[6]['f'] = Bishop(Black)
	game.board.matrix[6]['g'] = Bishop(White)
	game.board.matrix[6]['h'] = Empty()

	game.board.matrix[7]['a'] = Empty()
	game.board.matrix[7]['b'] = Empty()
	game.board.matrix[7]['c'] = Empty()
	game.board.matrix[7]['d'] = Empty()
	game.board.matrix[7]['e'] = Empty()
	game.board.matrix[7]['f'] = Empty()
	game.board.matrix[7]['g'] = Empty()
	game.board.matrix[7]['h'] = Empty()

	game.Current.King = "g1"
	game.Current.Next.King = "a6"

	game.EndTurn()

	go gameManager.Process(Message{
		Type: Move,
		Payload: MovePiece{
			To:     "a5",
			From:   "a6",
			GameId: game.Id.String(),
		},
	})

	<-p1.Outgoing

	if !reflect.DeepEqual(game.board.Square("a5"), King(Black)) {
		t.Errorf("Expected king on a5, got %v", game.board.Square("a5"))
	}
}

func TestNoMateButShouldBeMate(t *testing.T) {
	gameManager := NewGameManager()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go gameManager.Process(Message{
		Type: CreateGame,
		Payload: MatchParams{
			Players: []*Player{p1, p2},
			TimeControl: TimeControl{
				Duration:  "5m",
				Increment: "1s",
			},
		},
	})

	res := <-p1.Outgoing
	<-p2.Outgoing

	params := res.Payload.(GameStart)
	game := gameManager.FindGame(params.GameId)

	game.board.matrix[0]['a'] = Empty()
	game.board.matrix[0]['b'] = Empty()
	game.board.matrix[0]['c'] = Empty()
	game.board.matrix[0]['d'] = Empty()
	game.board.matrix[0]['e'] = Empty()
	game.board.matrix[0]['f'] = Empty()
	game.board.matrix[0]['g'] = Empty()
	game.board.matrix[0]['h'] = Empty()

	game.board.matrix[1]['a'] = Empty()
	game.board.matrix[1]['b'] = Rook(Black)
	game.board.matrix[1]['c'] = Empty()
	game.board.matrix[1]['d'] = Empty()
	game.board.matrix[1]['e'] = Empty()
	game.board.matrix[1]['f'] = Empty()
	game.board.matrix[1]['g'] = Empty()
	game.board.matrix[1]['h'] = Empty()

	game.board.matrix[2]['b'] = Pawn(Black)
	game.board.matrix[2]['g'] = King(White)
	game.board.matrix[2]['h'] = Pawn(White)

	game.board.matrix[3]['g'] = Pawn(White)

	game.board.matrix[6]['a'] = Rook(White)
	game.board.matrix[6]['b'] = Empty()
	game.board.matrix[6]['c'] = Empty()
	game.board.matrix[6]['d'] = Empty()
	game.board.matrix[6]['e'] = Empty()
	game.board.matrix[6]['f'] = Empty()
	game.board.matrix[6]['g'] = Empty()
	game.board.matrix[6]['h'] = Rook(White)

	game.board.matrix[7]['a'] = Empty()
	game.board.matrix[7]['b'] = Empty()
	game.board.matrix[7]['c'] = Empty()
	game.board.matrix[7]['d'] = Empty()
	game.board.matrix[7]['e'] = Empty()
	game.board.matrix[7]['f'] = King(Black)
	game.board.matrix[7]['g'] = Empty()
	game.board.matrix[7]['h'] = Empty()

	game.Current.King = "g3"
	game.Current.Next.King = "f8"

	if game.board.Square("h7") != Rook(White) {
		t.Errorf("Expected rook on h7, got %v", game.board.Square("h7"))
	}

	go gameManager.Process(Message{
		Type: Move,
		Payload: MovePiece{
			To:     "h8",
			From:   "h7",
			GameId: game.Id.String(),
		},
	})

	winner := <-p1.Outgoing
	if winner.Type != GameOver {
		t.Errorf("Expected game over, got %v", winner.Type)
	}
	payload := winner.Payload.(GameOverResponse)
	if !payload.Winner {
		t.Errorf("Expected white to win by checkmate")
	}

	loser := <-p2.Outgoing
	if loser.Type != GameOver {
		t.Errorf("Expected game over, got %v", loser.Type)
	}

	payload = loser.Payload.(GameOverResponse)
	if payload.Winner {
		t.Errorf("Expected black to lose by checkmate")
	}
}
