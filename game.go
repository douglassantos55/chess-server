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
	Next   *GamePlayer
	Player *Player
	Color  Color

	start time.Time
	left  time.Duration
	timer *time.Timer
	mutex *sync.Mutex
}

func NewGamePlayer(color Color, player *Player, duration time.Duration) *GamePlayer {
	timer := time.NewTimer(duration)
	timer.Stop()

	return &GamePlayer{
		Player: player,
		Color:  color,

		mutex: new(sync.Mutex),
		left:  duration,
		timer: timer,
	}
}

func (p *GamePlayer) SetNext(player *GamePlayer) {
	if p.Next != nil {
		player.Next = p.Next
	} else {
		player.Next = p
	}

	p.Next = player
}

func (p *GamePlayer) StartTimer() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.start = time.Now()
	p.timer.Reset(p.left)
}

func (p *GamePlayer) StopTimer() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.timer.Stop()
	p.left = time.Since(p.start)
}

func (p *GamePlayer) Send(response Response) {
	p.Player.Send(response)
}

type Game struct {
	Id      uuid.UUID
	Current *GamePlayer
	Over    chan GameResult

	board *Board
	mutex *sync.Mutex
}

func NewGame(duration time.Duration, players []*Player) *Game {
	white := NewGamePlayer(White, players[0], duration)
	black := NewGamePlayer(Black, players[1], duration)

	white.SetNext(black)

	game := &Game{
		Id:      uuid.New(),
		Over:    make(chan GameResult),
		Current: white,

		board: NewBoard(),
		mutex: new(sync.Mutex),
	}

	go func() {
		select {
		case <-white.timer.C:
			game.GameOver(white.Player)
		case <-black.timer.C:
			game.GameOver(black.Player)
		}
	}()

	return game
}

func (g *Game) GameOver(loser *Player) {
	winner := g.Current.Player

	if winner == loser {
		winner = g.Current.Next.Player
	}

	g.Over <- GameResult{
		Loser:  loser,
		Winner: winner,
	}
}

func (g *Game) EndTurn() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.Current.StopTimer()
	g.Current = g.Current.Next
}

func (g *Game) StartTurn() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.Current.StartTimer()
}

func (g *Game) Move(from, to string) {
	piece := g.board.Square(from)

	if piece.Color == g.Current.Color {
		g.board.Move(from, to)

		g.EndTurn()
		g.StartTurn()
	}
}

// TODO: register game as a listener
func (g *Game) Start() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	players := make(map[*GamePlayer]*GamePlayer)

	for player := g.Current; player != nil; player = player.Next {
		// we've seen this dude, it's looping now, get out
		if _, ok := players[player]; ok {
			break
		}

		player.Send(Response{
			Type: StartGame,
			Payload: GameParams{
				GameId: g.Id,
				Color:  player.Color,
			},
		})

		players[player] = player
	}

	g.Current.StartTimer()
}
