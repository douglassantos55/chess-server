package main

import "github.com/google/uuid"

type MessageType string
type ResponseType string

const (
	QueueUp        MessageType = "queue_up"
	Dequeue        MessageType = "dequeue"
	Disconnected   MessageType = "disconnected"
	MatchConfirmed MessageType = "match_confirmed"
	MatchDeclined  MessageType = "match_declined"
	CreateGame     MessageType = "create_game"
	MatchFound     MessageType = "match_found"
)

const (
	WaitForMatch     ResponseType = "wait_for_match"
	ConfirmMatch     ResponseType = "confirm_match"
	WaitOtherPlayers ResponseType = "wait_other_players"
	MatchCanceled    ResponseType = "match_canceled"
	StartGame        ResponseType = "start_game"
	StartTurn        ResponseType = "start_turn"
	GameOver         ResponseType = "game_over"
)

type Message struct {
	Type    MessageType
	Text    string
	Player  *Player
	Payload interface{}
}

type Response struct {
	Type    ResponseType
	Text    string
	Payload interface{}
}

type GameParams struct {
	GameId uuid.UUID
	Color  Color
}
