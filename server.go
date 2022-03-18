package main

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler interface {
	Process(event Message)
}

type Server struct {
	server   *http.Server
	handlers []Handler
}

func NewServer(handlers []Handler) *Server {
	return &Server{
		handlers: handlers,
		server:   &http.Server{},
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
	socket, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	player := NewPlayer(socket)

	for {
		defer player.Close()

		message := <-player.Incoming
		message.Player = player

		for _, handler := range s.handlers {
			go handler.Process(message)
		}
	}
}
