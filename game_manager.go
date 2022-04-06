package main

import (
	"sync"
	"time"

	"github.com/google/uuid"
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

func (g *GameManager) CreateGame(players []*Player) *Game {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	game := NewGame(time.Second, players)
	g.games[game.Id] = game

	go func() {
		result := <-game.Over

		g.RemoveGame(game.Id)

		result.Winner.Send(Response{
			Type:    GameOver,
			Payload: result,
		})

		result.Loser.Send(Response{
			Type:    GameOver,
			Payload: result,
		})
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
		game := g.CreateGame(event.Payload.([]*Player))
		game.Start()
	case Move:
		data := event.Payload.(MovePiece)
		game := g.FindGame(data.GameId)

		game.Move(data.From, data.To)
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
