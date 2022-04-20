package pkg

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
	case <-time.After(200 * time.Millisecond):
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

	<-player.Outgoing

	if queueManager.queue.Pop() != player {
		t.Error("Expected head to be player")
	}
}

func TestCancel(t *testing.T) {
	p1 := NewTestPlayer()
	queueManager := NewQueueManager()

	go queueManager.Process(Message{
		Type:   QueueUp,
		Player: p1,
	})

	<-p1.Outgoing

	<-wait(func() {
		queueManager.Process(Message{
			Type:   Dequeue,
			Player: p1,
		})
	})

	queue := queueManager.queue

	if queue.Length() != 0 {
		t.Errorf("Expected empty queue, got %v", queue.Length())
	}
}

func TestDisconnectRemovesFromQueue(t *testing.T) {
	player := NewTestPlayer()
	queueManager := NewQueueManager()

	go queueManager.Process(Message{
		Type:   QueueUp,
		Player: player,
	})

	<-player.Outgoing

	<-wait(func() {
		queueManager.Process(Message{
			Type:   Disconnected,
			Player: player,
		})
	})

	got := queueManager.queue.Pop()

	if got != nil {
		t.Errorf("Expected empty queue, got %v", got)
	}
}

func TestDispatchesMatchFound(t *testing.T) {
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
	case res := <-p2.Outgoing:
		if res.Type != WaitForMatch {
			t.Errorf("Expected wait for match, got %+v", res)
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Timeout before server response")
	}

	select {
	case res := <-Dispatcher:
		if res.Type != MatchFound {
			t.Errorf("Expected match found, got %+v", res)
		}

		players := res.Payload.([]*Player)

		if len(players) != 2 {
			t.Errorf("Expected 2 players, got %v", len(players))
		}
	case <-time.After(time.Second):
		t.Error("Timeout before server response")
	}

	if queueManager.queue.Length() != 0 {
		t.Errorf("Expected empty queue, got %v", queueManager.queue.Length())
	}
}
