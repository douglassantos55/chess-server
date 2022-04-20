package pkg

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	Dispatcher = make(chan Message)
	os.Exit(m.Run())
}

func StartServer(handlers []Handler) *Server {
	server := NewServer(handlers)
	go server.Listen("0.0.0.0:8080")

	// ...
	time.Sleep(time.Millisecond)

	return server
}

func TestAcceptsConnection(t *testing.T) {
	server := StartServer([]Handler{})
	defer server.Shutdown()

	_, err := NewClient()

	if err != nil {
		t.Errorf("Expected connection, got error %+v", err)
	}
}

func TestShutdown(t *testing.T) {
	server := StartServer([]Handler{})
	server.Shutdown()

	_, err := NewClient()

	if err == nil {
		t.Error("Expected shutdown")
	}
}

func TestHandleConnection(t *testing.T) {
	server := StartServer([]Handler{})
	defer server.Shutdown()

	_, err := NewClient()

	if err != nil {
		t.Errorf("Expected connection, got error: %+v", err)
	}
}
