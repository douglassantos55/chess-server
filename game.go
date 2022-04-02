package main

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Color string

const (
	Black Color = "black"
	White Color = "white"
)

type GameResult struct {
	Winner *Player
	Loser  *Player
}

type GamePlayer struct {
	Player *Player
	Color  Color

	timer *time.Timer
	left  time.Duration
	mutex *sync.Mutex
	start time.Time
}

func NewGamePlayer(color Color, player *Player, duration time.Duration) *GamePlayer {
	timer := time.NewTimer(duration)
	timer.Stop()

	return &GamePlayer{
		Player: player,
		Color:  Black,

		mutex: new(sync.Mutex),
		left:  duration,
		timer: timer,
	}
}

func (p *GamePlayer) StartTurn() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.start = time.Now()
	p.timer.Reset(p.left)
}

func (p *GamePlayer) EndTurn() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.timer.Stop()
	p.left = time.Since(p.start)
}

func (p *GamePlayer) Send(response Response) {
	p.Player.Send(response)
}

type Game struct {
	Id    uuid.UUID
	Over  chan GameResult
	Black *GamePlayer
	White *GamePlayer

	board *Board
	mutex *sync.Mutex
}

func NewGame(players []*Player) *Game {
	return &Game{
		Id:    uuid.New(),
		Over:  make(chan GameResult),
		White: NewGamePlayer(White, players[0], time.Second),
		Black: NewGamePlayer(Black, players[1], time.Second),

		board: NewBoard(),
		mutex: new(sync.Mutex),
	}
}

func (g *Game) GameOver(winner, loser *GamePlayer) {
	g.Over <- GameResult{
		Winner: winner.Player,
		Loser:  loser.Player,
	}
}

func (g *Game) StartTurn(color Color) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if color == White {
		g.White.StartTurn()
	} else {
		g.Black.StartTurn()
	}
}

func (g *Game) EndTurn(color Color) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if color == White {
		g.White.EndTurn()
	} else {
		g.Black.EndTurn()
	}
}

// TODO: register game as a listener
func (g *Game) Start() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

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
				g.GameOver(g.Black, g.White)
			case <-g.Black.timer.C:
				g.GameOver(g.White, g.Black)
			}
		}
	}()
}
