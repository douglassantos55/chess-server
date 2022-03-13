package main

import (
	"testing"
	"time"
)

func TestReturnsResponse(t *testing.T) {
	queueManager := NewQueueManager()
	player := NewPlayer()

	// channels without goroutines are a nono
	go queueManager.Process(Message{
		Type:   QueueUp,
		Player: player,
	})

	select {
	case res := <-player.Incoming:
		if res.Type != WaitForMatch {
			t.Errorf("Expected wait for match, got %+v", res)
		}
	case <-time.After(time.Second):
		t.Error("Timeout before server response")
	}
}

func TestInvalidType(t *testing.T) {
	queueManager := NewQueueManager()
	player := NewPlayer()

	go queueManager.Process(Message{
		Type:   "something",
		Player: player,
	})

	select {
	case <-time.After(time.Second):
	case <-player.Incoming:
		t.Error("Should not get response")
	}
}

func TestAddsToQueue(t *testing.T) {
	queueManager := NewQueueManager()
	player := NewPlayer()

	go queueManager.Process(Message{
		Type:   QueueUp,
		Player: player,
	})

	time.Sleep(time.Millisecond)

	if queueManager.queue.Pop() == nil {
		t.Error("Expected head to be player")
	}
}
