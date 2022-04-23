package pkg

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

func (m *MatchMaker) CreateMatch(players []*Player, timeControl TimeControl) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	match := NewMatch(players, timeControl)
	m.matches[match.Id] = match

	go match.AskConfirmation()
	go match.WaitConfirmation(m.timeout)

	go func() {
		select {
		case players := <-match.Ready:
			m.RemoveMatch(match.Id)

			Dispatcher <- Message{
				Type: CreateGame,
				Payload: MatchParams{
					Players:     players,
					TimeControl: match.TimeControl,
				},
			}
		case requeue := <-match.Canceled:
			m.RemoveMatch(match.Id)

			for _, player := range match.Players {
				player.Send(Response{
					Type: MatchCanceled,
				})
			}

			for _, player := range requeue {
				Dispatcher <- Message{
					Type:   QueueUp,
					Player: player,
					Payload: map[string]interface{}{
						"duration":  match.TimeControl.Duration,
						"increment": match.TimeControl.Increment,
					},
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

func (m *MatchMaker) CancelPlayerMatches(player *Player) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, match := range m.matches {
		for _, p := range match.Players {
			if p == player {
				match.Cancel()
			}
		}
	}
}

func (m *MatchMaker) Process(event Message) {
	switch event.Type {
	case MatchFound:
		params := event.Payload.(MatchParams)
		m.CreateMatch(params.Players, params.TimeControl)

	case MatchConfirmed:
		matchId, _ := uuid.Parse(event.Payload.(string))
		m.ConfirmMatch(matchId, event.Player)

	case MatchDeclined:
		matchId, _ := uuid.Parse(event.Payload.(string))
		m.CancelMatch(matchId)

	case Disconnected:
		m.CancelPlayerMatches(event.Player)
	}
}
