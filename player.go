package main

import "github.com/google/uuid"

type Player struct {
	Id       uuid.UUID
	Incoming chan Response
}

func NewPlayer() *Player {
	// &Player{} returns the same pointer
	// no matter how many times you call it.
	// Adding a unique ID creates different
	// instances/pointers
	return &Player{
		Id:       uuid.New(),
		Incoming: make(chan Response),
	}
}

func (p *Player) Send(response Response) {
	p.Incoming <- response
}
