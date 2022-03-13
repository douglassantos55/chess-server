package main

type node struct {
	Player *Player
	next   *node
}

type Queue struct {
	head *node
	tail *node
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Push(player *Player) {
	node := &node{Player: player}

	if q.head == nil {
		q.head = node
		q.tail = node
	} else {
		q.tail.next = node
		q.tail = node
	}
}

func (q *Queue) Pop() *Player {
	if q.head == nil {
		return nil
	}

	if q.head == q.tail {
		q.tail = q.head.next
	}

	player := q.head.Player
	q.head = q.head.next

	return player
}
