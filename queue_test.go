package main

import (
	"testing"
)

func NewTestPlayer() *Player {
	return &Player{
		Incoming: make(chan Message),
		Outgoing: make(chan Response),
	}
}

func TestEmpty(t *testing.T) {
	queue := NewQueue()
	player := NewTestPlayer()

	queue.Push(player)

	if queue.head == nil {
		t.Errorf("Expected something in head")
	}
	if queue.head != queue.tail {
		t.Error("Expected head and tail to point to the same thing")
	}
	if queue.head.Player != player {
		t.Error("Expected head to have player")
	}
}

func TestAppends(t *testing.T) {
	queue := NewQueue()
	player := NewTestPlayer()

	queue.Push(player)
	queue.Push(NewTestPlayer())
	queue.Push(NewTestPlayer())

	if queue.head.Player != player {
		t.Error("Expected head to have player")
	}
	if queue.head.next == nil {
		t.Error("Expected more than one item")
	}
	if queue.tail == nil {
		t.Error("Expected tail")
	}
}

func TestPop(t *testing.T) {
	queue := NewQueue()
	player := NewTestPlayer()

	queue.Push(player)
	queue.Push(NewTestPlayer())

	if queue.Pop() != player {
		t.Error("Expected pop to return player")
	}
	if queue.tail.next != nil {
		t.Error("Expected no next in tail")
	}
	if queue.head != queue.tail {
		t.Error("Expected head and tail to point to the same thing")
	}
}

func TestPopEmpty(t *testing.T) {
	queue := NewQueue()

	if queue.Pop() != nil {
		t.Error("Should not be able to pop from empty queue")
	}
}

func TestPopUnique(t *testing.T) {
	queue := NewQueue()
	queue.Push(NewTestPlayer())

	queue.Pop()

	if queue.head != nil {
		t.Error("Expected nil head")
	}
	if queue.tail != queue.head {
		t.Error("Expected head and tail to be nil")
	}
}
