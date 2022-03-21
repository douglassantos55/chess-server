package main

import (
	"testing"
	"time"
)

type TestHandler struct {
	QueueUp      chan bool
	Disconnected chan bool
}

func (h *TestHandler) Process(event Message) {
	switch event.Type {
	case QueueUp:
		h.QueueUp <- true
	case Disconnected:
		h.Disconnected <- true
	}
}

func TestHandlers(t *testing.T) {
	testHandler := &TestHandler{
		QueueUp:      make(chan bool),
		Disconnected: make(chan bool),
	}

	server := StartServer([]Handler{
		testHandler,
	})

	defer server.Shutdown()

	client, _ := NewClient()
	client.Send(QueueUp)

	select {
	case invoked := <-testHandler.QueueUp:
		if !invoked {
			t.Error("Expected queue up, got false instead")
		}
	case <-time.After(time.Second):
		t.Error("Expected handler to execute, got timeout instead")
	}
}

func TestClientDisconnect(t *testing.T) {
	testHandler := &TestHandler{
		QueueUp:      make(chan bool),
		Disconnected: make(chan bool),
	}

	server := StartServer([]Handler{
		testHandler,
	})

	defer server.Shutdown()

	client, _ := NewClient()
	client.Close()

	select {
	case disconnected := <-testHandler.Disconnected:
		if !disconnected {
			t.Error("Expected disconnected, got false instead")
		}
	case <-time.After(time.Second):
		t.Error("Expected handler to execute, got timeout instead")
	}

	assertPanic(t, func() { client.Send(QueueUp) })
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected send on closed client to panic")
		}
	}()

	f()
}
