package main

import (
	"testing"
	"time"
)

type TestHandler struct {
	Invoked chan bool
}

func (h *TestHandler) Process(event Message) {
	h.Invoked <- true
}

func TestHandlers(t *testing.T) {
	testHandler := &TestHandler{
		Invoked: make(chan bool),
	}

	server := StartServer([]Handler{
		testHandler,
	})

	defer server.Shutdown()

	client, _ := NewClient()
	client.Send(QueueUp)

	select {
	case invoked := <-testHandler.Invoked:
		if !invoked {
			t.Error("Expected handler to execute, got false instead")
		}
	case <-time.After(time.Second):
		t.Error("Expected handler to execute, got timeout instead")
	}
}
