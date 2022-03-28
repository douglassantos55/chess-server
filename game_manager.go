package main

import (
	"sync"

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

func (g *GameManager) CreateGame(players []*Player) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	game := NewGame(players)
	g.games[game.Id] = game

	game.White.Send(Response{
		Type: StartGame,
		Payload: GameParams{
			GameId: game.Id,
			Color:  White,
		},
	})

	game.Black.Send(Response{
		Type: StartGame,
		Payload: GameParams{
			GameId: game.Id,
			Color:  Black,
		},
	})

}

func (g *GameManager) FindGame(gameId uuid.UUID) *Game {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	return g.games[gameId]
}

func (g *GameManager) Process(event Message) {
	switch event.Type {
	case CreateGame:
		g.CreateGame(event.Payload.([]*Player))
	}
}
