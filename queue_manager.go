package main

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
	}
}
