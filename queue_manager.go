package main

const MAX_PLAYERS = 2

type QueueManager struct {
	queue *Queue
}

func NewQueueManager() *QueueManager {
	return &QueueManager{
		queue: NewQueue(),
	}
}

func (q *QueueManager) Process(event Message) {
	switch event.Type {
	case QueueUp:
		q.queue.Push(event.Player)

		event.Player.Send(Response{
			Type: WaitForMatch,
			Text: "Wait for match",
		})

		if q.queue.Length() == MAX_PLAYERS {
			for i := 0; i < MAX_PLAYERS; i++ {
				player := q.queue.Pop()

				player.Send(Response{
					Type: MatchFound,
				})
			}

			Dispatcher <- Message{
				Type: MatchFound,
			}
		}
	case Dequeue, Disconnected:
		q.queue.Remove(event.Player)

	}
}
