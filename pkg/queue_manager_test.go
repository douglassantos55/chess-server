package pkg

import (
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"
)

func TestReturnsResponse(t *testing.T) {
	player := NewTestPlayer()
	queueManager := NewQueueManager()

	go queueManager.Process(Message{
		Type:   QueueUp,
		Player: player,
		Payload: map[string]interface{}{
			"duration":  "1m",
			"increment": "0s",
		},
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
	player1 := NewTestPlayer()
	player2 := NewTestPlayer()

	queueManager := NewQueueManager()

	payload1 := map[string]interface{}{
		"duration":  "1m",
		"increment": "0s",
	}

	payload2 := map[string]interface{}{
		"duration":  "10m",
		"increment": "0s",
	}

	go queueManager.Process(Message{
		Type:    QueueUp,
		Player:  player1,
		Payload: payload1,
	})

	<-player1.Outgoing

	go queueManager.Process(Message{
		Type:    QueueUp,
		Player:  player2,
		Payload: payload2,
	})

	<-player2.Outgoing

	var params1 QueueUpParams
	mapstructure.Decode(payload1, &params1)

	if queueManager.queue[params1].Length() != 1 {
		t.Error("Expected 1 player on 1m+0s queue")
	}

	var params2 QueueUpParams
	mapstructure.Decode(payload2, &params2)

	if queueManager.queue[params2].Length() != 1 {
		t.Error("Expected 1 player on 10m+0s queue")
	}
}

func TestCancel(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	queueManager := NewQueueManager()

	payload1 := map[string]interface{}{
		"duration":  "1m",
		"increment": "0s",
	}

	go queueManager.Process(Message{
		Type:    QueueUp,
		Player:  p1,
		Payload: payload1,
	})

	<-p1.Outgoing

	payload2 := map[string]interface{}{
		"duration":  "10m",
		"increment": "0s",
	}

	go queueManager.Process(Message{
		Type:    QueueUp,
		Player:  p2,
		Payload: payload2,
	})

	<-p2.Outgoing

	<-wait(func() {
		queueManager.Process(Message{
			Type:   Dequeue,
			Player: p1,
		})
	})

	queue := queueManager.queue

	var params1 QueueUpParams
	mapstructure.Decode(payload1, &params1)

	if queue[params1].Length() != 0 {
		t.Errorf("Expected empty queue, got %v", queue[params1].Length())
	}

	var params2 QueueUpParams
	mapstructure.Decode(payload2, &params2)

	if queue[params2].Length() == 0 {
		t.Errorf("Expected 1 player in queue, got %v", queue[params2].Length())
	}
}

func TestDisconnectRemovesFromQueue(t *testing.T) {
	player := NewTestPlayer()
	queueManager := NewQueueManager()

	payload := map[string]interface{}{
		"duration":  "1m",
		"increment": "0s",
	}

	go queueManager.Process(Message{
		Type:    QueueUp,
		Player:  player,
		Payload: payload,
	})

	<-player.Outgoing

	<-wait(func() {
		queueManager.Process(Message{
			Type:   Disconnected,
			Player: player,
		})
	})

	var params QueueUpParams
	mapstructure.Decode(payload, &params)

	got := queueManager.queue[params].Pop()

	if got != nil {
		t.Errorf("Expected empty queue, got %v", got)
	}
}

func TestDispatchesMatchFound(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	queueManager := NewQueueManager()

	payload := map[string]interface{}{
		"duration":  "1m",
		"increment": "0s",
	}

	go queueManager.Process(Message{
		Type:    QueueUp,
		Player:  p1,
		Payload: payload,
	})

	select {
	case res := <-p1.Outgoing:
		if res.Type != WaitForMatch {
			t.Errorf("Expected wait for match, got %+v", res)
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Timeout before server response")
	}

	go queueManager.Process(Message{
		Type:    QueueUp,
		Player:  p2,
		Payload: payload,
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

	var params QueueUpParams
	mapstructure.Decode(payload, &params)

	if queueManager.queue[params].Length() != 0 {
		t.Errorf("Expected empty queue, got %v", queueManager.queue[params].Length())
	}
}

func TestDoesNotDispatchMatchFound(t *testing.T) {
	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	queueManager := NewQueueManager()

	go queueManager.Process(Message{
		Type:   QueueUp,
		Player: p1,
		Payload: map[string]interface{}{
			"duration":  "1m",
			"increment": "0s",
		},
	})

	select {
	case res := <-p1.Outgoing:
		if res.Type != WaitForMatch {
			t.Errorf("Expected wait for match, got %+v", res)
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Timeout before server response")
	}

	go queueManager.Process(Message{
		Type:   QueueUp,
		Player: p2,
		Payload: map[string]interface{}{
			"duration":  "5m",
			"increment": "0s",
		},
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
	case <-time.After(time.Second):
	case res := <-Dispatcher:
		if res.Type == MatchFound {
			t.Error("Should not receive match found")
		}
	}
}
