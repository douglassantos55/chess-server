package main

import (
	"testing"
	"time"
)

type TestQueueManager struct {
	Manager *QueueManager
}

func NewTestQueueManager() *TestQueueManager {
	return &TestQueueManager{
		Manager: NewQueueManager(),
	}
}

func (q *TestQueueManager) Process(event Message) chan bool {
	channel := make(chan bool)

	go func() {
		q.Manager.Process(event)
		channel <- true
	}()

	return channel
}

func TestReturnsResponse(t *testing.T) {
	player := NewTestPlayer()
	queueManager := NewTestQueueManager()

	queueManager.Process(Message{
		Type:   QueueUp,
		Player: player,
	})

	select {
	case res := <-player.Outgoing:
		if res.Type != WaitForMatch {
			t.Errorf("Expected wait for match, got %+v", res)
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Timeout before server response")
	}
}

func TestInvalidType(t *testing.T) {
	player := NewTestPlayer()
	queueManager := NewTestQueueManager()

	queueManager.Process(Message{
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
	queueManager := NewTestQueueManager()

	queueManager.Process(Message{
		Type:   QueueUp,
		Player: player,
	})

	<-player.Outgoing

	if queueManager.Manager.queue.Pop() != player {
		t.Error("Expected head to be player")
	}
}

func TestCancel(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	queueManager := NewTestQueueManager()

	p1Ready := queueManager.Process(Message{
		Type:   QueueUp,
		Player: p1,
	})
	p2Ready := queueManager.Process(Message{
		Type:   QueueUp,
		Player: p2,
	})

	select {
	case <-p2.Outgoing:
		<-p2Ready
	case <-time.After(time.Second):
		t.Error("Should not timeout")
	}

	select {
	case <-p1.Outgoing:
		<-p1Ready
	case <-time.After(time.Second):
		t.Error("Should not timeout")
	}

	<-queueManager.Process(Message{
		Type:   Dequeue,
		Player: p2,
	})

	if queueManager.Manager.queue.tail.Player != p1 {
		t.Error("Expected tail to point to p1")
	}
}

func TestDisconnectRemovesFromQueue(t *testing.T) {
	player := NewTestPlayer()
	queueManager := NewTestQueueManager()

	ready := queueManager.Process(Message{
		Type:   QueueUp,
		Player: player,
	})

	select {
	case <-player.Outgoing:
		<-ready
	case <-time.After(time.Second):
		t.Error("Should not timeout")
	}

	<-queueManager.Process(Message{
		Type:   Disconnected,
		Player: player,
	})

	got := queueManager.Manager.queue.Pop()

	if got != nil {
		t.Errorf("Expected empty queue, got %v", got)
	}
}
