package main

import (
	"sync"

	"github.com/google/uuid"
)

type MatchMaker struct {
	mutex   *sync.Mutex
	matches map[uuid.UUID][]*Player
}

func NewMatchMaker() *MatchMaker {
	return &MatchMaker{
		mutex:   new(sync.Mutex),
		matches: make(map[uuid.UUID][]*Player),
	}
}

func (m *MatchMaker) HasMatches() bool {
	return len(m.matches) > 0
}

func (m *MatchMaker) CreateMatch(players []*Player) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	matchId := uuid.New()
	m.matches[matchId] = players

	for _, player := range players {
		player.Send(Response{
			Type:    ConfirmMatch,
			Payload: matchId,
		})
	}
}

func (m *MatchMaker) Process(event Message) {
	switch event.Type {
	case MatchFound:
		players := event.Payload.([]*Player)
		m.CreateMatch(players)
	}
}
