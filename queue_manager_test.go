package main

import (
	"testing"
	"time"
)

func TestReturnsResponse(t *testing.T) {
	player := NewTestPlayer()
	queueManager := NewQueueManager()

	go queueManager.Process(Message{
		Type:   QueueUp,
		Player: player,
	})

	select {
	case res := <-player.Outgoing:
		if res.Type != WaitForMatch {
			t.Errorf("Expected wait for match, got %+v", res)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timeout before server response")
	}
}

func TestInvalidType(t *testing.T) {
	player := NewTestPlayer()
	queueManager := NewQueueManager()

	go queueManager.Process(Message{
		Type:   "something",
		Player: player,
	})

	select {
	case <-time.After(100 * time.Millisecond):
	case <-player.Outgoing:
		t.Error("Should not get response")
	}
}

func TestAddsToQueue(t *testing.T) {
	player := NewTestPlayer()
	queueManager := NewQueueManager()

	go queueManager.Process(Message{
		Type:   QueueUp,
		Player: player,
	})

	time.Sleep(time.Millisecond)

	if queueManager.queue.Pop() != player {
		t.Error("Expected head to be player")
	}
}

func TestCancel(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	queueManager := NewQueueManager()

	go queueManager.Process(Message{
		Type:   QueueUp,
		Player: p1,
	})
	go queueManager.Process(Message{
		Type:   QueueUp,
		Player: p2,
	})

	select {
	case <-p2.Outgoing:
	case <-time.After(time.Second):
		t.Error("Should not timeout")
	}

	select {
	case <-p1.Outgoing:
	case <-time.After(time.Second):
		t.Error("Should not timeout")
	}

	go queueManager.Process(Message{
		Type:   Dequeue,
		Player: p2,
	})

	time.Sleep(time.Millisecond)

	if queueManager.queue.Tail() != p1 {
		t.Error("Expected tail to point to p1")
	}
}
