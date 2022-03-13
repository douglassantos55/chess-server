package main

import "github.com/google/uuid"

type Player struct {
	Id uuid.UUID
}

func NewPlayer() *Player {
	// &Player{} returns the same pointer
	// no matter how many times you call it.
	// Adding a unique ID creates different
	// instances/pointers
	return &Player{
		Id: uuid.New(),
	}
}
