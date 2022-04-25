package pkg

import (
	"sync"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

type GameManager struct {
	games map[uuid.UUID]*Game
	mutex *sync.Mutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		mutex: new(sync.Mutex),
		games: make(map[uuid.UUID]*Game),
	}
}

func (g *GameManager) CreateGame(players []*Player, timeControl TimeControl) *Game {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	game := NewGame(players, timeControl)
	g.games[game.Id] = game

	go func() {
		result := <-game.Over

		g.RemoveGame(game.Id)

		if result.Winner != nil {
			result.Winner.Send(Response{
				Type: GameOver,
				Payload: GameOverResponse{
					Reason: result.Reason,
					Winner: true,
					GameId: game.Id,
				},
			})
		}

		if result.Loser != nil {
			result.Loser.Send(Response{
				Type: GameOver,
				Payload: GameOverResponse{
					Reason: result.Reason,
					Winner: false,
					GameId: game.Id,
				},
			})
		}
	}()

	return game
}

func (g *GameManager) RemoveGame(gameId uuid.UUID) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	delete(g.games, gameId)
}

func (g *GameManager) FindGame(gameId uuid.UUID) *Game {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	return g.games[gameId]
}

func (g *GameManager) FindPlayerGame(player *Player) *Game {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	for _, game := range g.games {
		for p := game.Current; p != nil; p = p.Next {
			if p.Player == player {
				return game
			}
		}
	}

	return nil
}

func (g *GameManager) Process(event Message) {
	switch event.Type {
	case CreateGame:
		payload := event.Payload.(MatchParams)
		game := g.CreateGame(payload.Players, payload.TimeControl)
		game.Start()
	case Move:
		var data MovePiece
		mapstructure.Decode(event.Payload, &data)

		gameUuid, err := uuid.Parse(data.GameId)
		if err != nil {
			return
		}

		game := g.FindGame(gameUuid)
		if game == nil {
			return
		}

		moves := game.Move(data.From, data.To)
		if len(moves) > 0 {
			game.EndTurn()

			if game.IsCheckmate() {
				game.Checkmate()
			} else {
				game.StartTurn()

				for _, move := range moves {
					game.Current.Send(Response{
						Type: StartTurn,
						Payload: MoveResponse{
							From:   move.From.String(),
							To:     move.To.String(),
							Time:   game.Current.left,
							GameId: gameUuid,
						},
					})
				}
			}
		}
	case Resign:
		gameId := event.Payload.(uuid.UUID)
		game := g.FindGame(gameId)

		game.GameOver(event.Player, "Resignation")
	case Disconnected:
		game := g.FindPlayerGame(event.Player)

		if game != nil {
			game.GameOver(event.Player, "Abandonment")
		}
	}
}
