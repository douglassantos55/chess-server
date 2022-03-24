package main

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type MatchMaker struct {
	mutex   *sync.Mutex
	timeout time.Duration
	matches map[uuid.UUID]*Match
}

func NewMatchMaker(timeout time.Duration) *MatchMaker {
	return &MatchMaker{
		timeout: timeout,
		mutex:   new(sync.Mutex),
		matches: make(map[uuid.UUID]*Match),
	}
}

func (m *MatchMaker) HasMatches() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return len(m.matches) > 0
}

func (m *MatchMaker) RemoveMatch(matchId uuid.UUID) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.matches, matchId)
}

func (m *MatchMaker) CreateMatch(players []*Player) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	match := NewMatch(players)
	m.matches[match.Id] = match

	go match.AskConfirmation()
	go match.WaitConfirmation(m.timeout)

	go func() {
		select {
		case players := <-match.Ready:
			Dispatcher <- Message{
				Type:    GameStart,
				Payload: players,
			}
		case requeue := <-match.Canceled:
			m.RemoveMatch(match.Id)

			for _, player := range requeue {
				player.Send(Response{
					Type: MatchCanceled,
				})

				Dispatcher <- Message{
					Type:   QueueUp,
					Player: player,
				}
			}
		}
	}()
}

func (m *MatchMaker) CancelMatch(matchId uuid.UUID) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	match := m.matches[matchId]
	match.Cancel()
}

func (m *MatchMaker) ConfirmMatch(matchId uuid.UUID, player *Player) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	match := m.matches[matchId]
	match.Confirm(player)
}

func (m *MatchMaker) Process(event Message) {
	switch event.Type {
	case MatchFound:
		players := event.Payload.([]*Player)
		m.CreateMatch(players)

	case MatchConfirmed:
		matchId := event.Payload.(uuid.UUID)
		m.ConfirmMatch(matchId, event.Player)

	case MatchDeclined:
		matchId := event.Payload.(uuid.UUID)
		m.CancelMatch(matchId)
	}
}
