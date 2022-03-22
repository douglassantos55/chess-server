package main

import (
	"sync"
)

type node struct {
	Player *Player
	next   *node
}

type Queue struct {
	length int
	head   *node
	tail   *node
	mut    *sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		mut: new(sync.Mutex),
	}
}

func (q *Queue) Push(player *Player) {
	q.mut.Lock()
	defer q.mut.Unlock()

	node := &node{Player: player}

	if q.head == nil {
		q.head = node
		q.tail = node
	} else {
		q.tail.next = node
		q.tail = node
	}

	q.length++
}

func (q *Queue) Pop() *Player {
	q.mut.Lock()
	defer q.mut.Unlock()

	if q.head == nil {
		return nil
	}

	if q.head == q.tail {
		q.tail = q.head.next
	}

	player := q.head.Player
	q.head = q.head.next

	q.length--

	return player
}

func (q *Queue) Remove(player *Player) {
	q.mut.Lock()
	defer q.mut.Unlock()

	var prev *node

	for cur := q.head; cur != nil; cur = cur.next {
		if cur.Player == player {
			if prev != nil {
				prev.next = cur.next

				if cur.next == nil {
					q.tail = prev
				}
			} else {
				q.head = cur.next
			}

			q.length--
			break
		}

		prev = cur
	}
}

func (q *Queue) Length() int {
	q.mut.Lock()
	defer q.mut.Unlock()

	return q.length
}
