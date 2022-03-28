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
	timer  *time.Timer
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
		White: &GamePlayer{
			Player: players[0],
			Color:  White,
			timer:  time.NewTimer(time.Second),
		},
		Black: &GamePlayer{
			Player: players[1],
			Color:  Black,
			timer:  time.NewTimer(time.Second),
		},
	}
}
