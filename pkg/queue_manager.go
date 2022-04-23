package pkg

import (
	"sync"

	"github.com/mitchellh/mapstructure"
)

const MAX_PLAYERS = 2

type QueueManager struct {
	mutex *sync.Mutex
	queue map[TimeControl]*Queue
}

func NewQueueManager() *QueueManager {
	return &QueueManager{
		mutex: new(sync.Mutex),
		queue: make(map[TimeControl]*Queue),
	}
}

func (q *QueueManager) GetQueue(event Message) (*Queue, TimeControl) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	var timeControl TimeControl
	mapstructure.Decode(event.Payload.(map[string]interface{}), &timeControl)

	queue := q.queue[timeControl]

	if queue == nil {
		queue = NewQueue()
		q.queue[timeControl] = queue
	}

	return queue, timeControl
}

func (q *QueueManager) Process(event Message) {
	switch event.Type {
	case QueueUp:
		queue, timeControl := q.GetQueue(event)
		queue.Push(event.Player)

		event.Player.Send(Response{
			Type: WaitForMatch,
			Text: "Wait for match",
		})

		if queue.Length() == MAX_PLAYERS {
			players := []*Player{}

			for i := 0; i < MAX_PLAYERS; i++ {
				players = append(players, queue.Pop())
			}

			Dispatcher <- Message{
				Type: MatchFound,
				Payload: MatchParams{
					Players:     players,
					TimeControl: timeControl,
				},
			}
		}
	case Dequeue, Disconnected:
		for _, queue := range q.queue {
			queue.Remove(event.Player)
		}
	}
}
