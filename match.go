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

	Ready  chan []*Player
	Cancel chan []*Player

	Confirmed chan *Player
}

func NewMatch(players []*Player) *Match {
	return &Match{
		mutex: new(sync.Mutex),

		Id:      uuid.New(),
		Players: players,

		Ready:  make(chan []*Player),
		Cancel: make(chan []*Player),

		Confirmed: make(chan *Player, MAX_PLAYERS),
	}
}

func (m *Match) Confirm(player *Player) {
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
	confirmed := []*Player{}

outer:
	for {
		select {
		case player := <-m.Confirmed:
			confirmed = append(confirmed, player)

			if len(confirmed) == MAX_PLAYERS {
				m.Ready <- confirmed
				break outer
			}
		case <-time.After(timeout):
			m.Cancel <- confirmed
			break outer
		}
	}
}
