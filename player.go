package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	Id uuid.UUID

	Incoming chan Message
	Outgoing chan Response

	socket *websocket.Conn
}

func NewPlayer(socket *websocket.Conn) *Player {
	// &Player{} returns the same pointer
	// no matter how many times you call it.
	// Adding a unique ID creates different
	// instances/pointers
	player := &Player{
		Id: uuid.New(),

		Incoming: make(chan Message),
		Outgoing: make(chan Response),

		socket: socket,
	}

	go player.Read()
	go player.Write()

	return player
}

func (p *Player) Close() {
	p.socket.Close()
}

func (p *Player) Send(response Response) {
	p.Outgoing <- response
}

func (p *Player) Read() {
	for {
		defer p.Close()

		var msg Message
		err := p.socket.ReadJSON(&msg)

		if err == nil {
			p.Incoming <- msg
		}
	}
}

func (p *Player) Write() {
	for {
		msg, ok := <-p.Outgoing

		if ok {
			p.socket.WriteJSON(msg)
		}
	}
}
