package pkg

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

var Dispatcher chan Message = make(chan Message)

type Handler interface {
	Process(event Message)
}

type Server struct {
	closed   chan bool
	reading  chan bool
	server   *http.Server
	handlers []Handler
}

func NewServer(handlers []Handler) *Server {
	return &Server{
		handlers: handlers,
		server:   &http.Server{},
		closed:   make(chan bool),
	}
}

func (s *Server) Shutdown() {
	s.server.Shutdown(context.Background())
}

func (s *Server) Listen(addr string) {
	s.server.Addr = addr
	s.server.Handler = http.HandlerFunc(s.HandleRequest)
	err := s.server.ListenAndServe()

	if err == http.ErrServerClosed {
		s.closed <- true
		s.reading <- false
	}
}

func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	socket, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	player := NewPlayer(socket)

	go func() {
		for {
			select {
			case <-s.reading:
				return
			case message, ok := <-player.Incoming:
				message.Player = player

				if !ok { // disconnected
					break
				}

				Dispatcher <- message
			}
		}
	}()

	go func() {
		for {
			select {
			case <-s.closed:
				return
			case event, ok := <-Dispatcher:
				if !ok {
					break
				}

				for _, handler := range s.handlers {
					go handler.Process(event)
				}
			}
		}
	}()
}
