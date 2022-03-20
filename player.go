package main

import (
	"time"

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
	p.socket.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(time.Second),
	)
}

func (p *Player) Send(response Response) {
	p.Outgoing <- response
}

// Read incoming client messages
func (p *Player) Read() {
	defer p.Close()

	for {
		var msg Message
		err := p.socket.ReadJSON(&msg)

		if err != nil {
			p.Incoming <- Message{
				Type: Disconnected,
			}
			break
		}
		p.Incoming <- msg
	}
}

// Write responses to client
func (p *Player) Write() {
	for {
		msg := <-p.Outgoing
		err := p.socket.WriteJSON(msg)

		if err != nil {
			p.Close()
		}
	}
}
