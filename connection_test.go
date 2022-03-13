package main

import (
	"testing"
	"time"
)

func StartServer() *Server {
	server := NewServer()
	go server.Listen("0.0.0.0:8080")

	// ...
	time.Sleep(time.Millisecond)

	return server
}

func TestAcceptsConnection(t *testing.T) {
	server := StartServer()
	defer server.Shutdown()

	_, err := NewClient()

	if err != nil {
		t.Errorf("Expected connection, got error %+v", err)
	}
}

func TestShutdown(t *testing.T) {
	server := StartServer()
	server.Shutdown()

	_, err := NewClient()

	if err == nil {
		t.Error("Expected shutdown")
	}
}

func TestHandleConnection(t *testing.T) {
	server := StartServer()
	defer server.Shutdown()

	_, err := NewClient()

	if err != nil {
		t.Errorf("Expected connection, got error: %+v", err)
	}
}
