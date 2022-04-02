package main

import (
	"time"

	"github.com/google/uuid"
)

type Color string

const (
	Black Color = "black"
	White Color = "white"
)

type GamePlayer struct {
	Player *Player
	Color  Color

	timer *time.Timer
	left  time.Duration
}

func NewGamePlayer(color Color, player *Player, duration time.Duration) *GamePlayer {
	return &GamePlayer{
		Player: player,
		Color:  Black,
		left:   duration,
		timer:  new(time.Timer),
	}
}

func (p *GamePlayer) StartTurn() {
	p.timer = time.NewTimer(p.left)
}

func (p *GamePlayer) Send(response Response) {
	p.Player.Send(response)
}

type Game struct {
	board *Board

	Id    uuid.UUID
	Black *GamePlayer
	White *GamePlayer
}

func NewGame(players []*Player) *Game {
	return &Game{
		Id:    uuid.New(),
		board: NewBoard(),
		White: NewGamePlayer(White, players[0], time.Second),
		Black: NewGamePlayer(Black, players[1], time.Second),
	}
}

func (g *Game) Start() {
	g.White.Send(Response{
		Type: StartGame,
		Payload: GameParams{
			GameId: g.Id,
			Color:  White,
		},
	})

	g.Black.Send(Response{
		Type: StartGame,
		Payload: GameParams{
			GameId: g.Id,
			Color:  Black,
		},
	})

	g.White.StartTurn()

	go func() {
		for {
			select {
			case <-g.White.timer.C:
				g.White.Send(Response{Type: GameOver})
				g.Black.Send(Response{Type: GameOver})
			case <-g.Black.timer.C:
				g.White.Send(Response{Type: GameOver})
				g.Black.Send(Response{Type: GameOver})
			}
		}
	}()
}
