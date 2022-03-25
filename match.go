package main

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Match struct {
	mutex *sync.Mutex

	Id      uuid.UUID
	Players []*Player

	Done     chan bool
	Ready    chan []*Player
	Canceled chan []*Player

	Confirmed chan *Player
}

func NewMatch(players []*Player) *Match {
	return &Match{
		mutex: new(sync.Mutex),

		Id:      uuid.New(),
		Players: players,

		Done:     make(chan bool),
		Ready:    make(chan []*Player),
		Canceled: make(chan []*Player),

		Confirmed: make(chan *Player, MAX_PLAYERS),
	}
}

func (m *Match) Cancel() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	close(m.Confirmed)
}

func (m *Match) Confirm(player *Player) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.Confirmed <- player

	player.Send(Response{
		Type: WaitOtherPlayers,
	})
}

func (m *Match) AskConfirmation() {
	for _, player := range m.Players {
		player.Send(Response{
			Type:    ConfirmMatch,
			Payload: m.Id,
		})
	}
}

func (m *Match) WaitConfirmation(timeout time.Duration) {
	go func() {
		select {
		case <-m.Done:
		case <-time.After(timeout):
			m.Cancel()
		}
	}()

	confirmed := []*Player{}

	for player := range m.Confirmed {
		confirmed = append(confirmed, player)

		if len(confirmed) == MAX_PLAYERS {
			m.Ready <- confirmed
		}
	}

	m.Canceled <- confirmed
	m.Done <- true
}
