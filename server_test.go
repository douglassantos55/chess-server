package main

import (
	"sync"
	"testing"
	"time"
)

type TestHandler struct {
	Invoked int
	mutex   *sync.Mutex
}

func (h *TestHandler) Process(event Message) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.Invoked++
}

func (h *TestHandler) Count() int {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return h.Invoked
}

func TestHandlers(t *testing.T) {
	testHandler := &TestHandler{
		mutex: new(sync.Mutex),
	}

	server := StartServer([]Handler{
		testHandler,
	})

	defer server.Shutdown()

	client, _ := NewClient()
	client.Send(QueueUp)

	time.Sleep(time.Millisecond)

	if testHandler.Count() != 1 {
		t.Error("Expected handler to execute")
	}
}
