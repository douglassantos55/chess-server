package main

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	server *http.Server
}

func NewServer() *Server {
	return &Server{
		server: &http.Server{},
	}
}

func (s *Server) Shutdown() {
	s.server.Shutdown(context.Background())
}

func (s *Server) Listen(addr string) {
	s.server.Addr = addr
	s.server.Handler = http.HandlerFunc(s.HandleRequest)
	s.server.ListenAndServe()
}

func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	_, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	for {
	}
}
