package pkg

import (
	"sync"
)

const MAX_PLAYERS = 2

type QueueManager struct {
	mutex *sync.Mutex
	queue map[QueueUpParams]*Queue
}

func NewQueueManager() *QueueManager {
	return &QueueManager{
		mutex: new(sync.Mutex),
		queue: make(map[QueueUpParams]*Queue),
	}
}

func (q *QueueManager) GetQueue(event Message) *Queue {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	params := event.Payload.(QueueUpParams)
	queue := q.queue[params]

	if queue == nil {
		queue = NewQueue()
		q.queue[params] = queue
	}

	return queue
}

func (q *QueueManager) Process(event Message) {
	switch event.Type {
	case QueueUp:
		queue := q.GetQueue(event)
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
				Type:    MatchFound,
				Payload: players,
			}
		}
	case Dequeue, Disconnected:
		for _, queue := range q.queue {
			queue.Remove(event.Player)
		}
	}
}
