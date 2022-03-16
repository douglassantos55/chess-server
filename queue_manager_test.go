package main

import (
	"testing"
	"time"
)

func TestReturnsResponse(t *testing.T) {
	server := StartServer([]Handler{
		NewQueueManager(),
	})
	defer server.Shutdown()

	client, _ := NewClient()
	client.Send(QueueUp)

	select {
	case res := <-client.Incoming:
		if res.Type != WaitForMatch {
			t.Errorf("Expected wait for match, got %+v", res)
		}
	case <-time.After(time.Second):
		t.Error("Timeout before server response")
	}
}

func TestInvalidType(t *testing.T) {
	server := StartServer([]Handler{
		NewQueueManager(),
	})
	defer server.Shutdown()

	client, _ := NewClient()
	client.Send("something")

	select {
	case <-time.After(100 * time.Millisecond):
	case <-client.Incoming:
		t.Error("Should not get response")
	}
}

func TestAddsToQueue(t *testing.T) {
	queueManager := NewQueueManager()

	server := StartServer([]Handler{
		queueManager,
	})
	defer server.Shutdown()

	client, _ := NewClient()
	client.Send(QueueUp)

    time.Sleep(time.Millisecond)

	if queueManager.queue.Pop() == nil {
		t.Error("Expected head to be player")
	}
}
