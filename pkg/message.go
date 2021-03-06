package pkg

import (
	"github.com/google/uuid"
)

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
	Move           MessageType = "move_piece"
	Resign         MessageType = "resign"
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
	Type    MessageType `json:"type"`
	Text    string      `json:"text"`
	Player  *Player     `json:"player"`
	Payload interface{} `json:"payload"`
}

type Response struct {
	Type    ResponseType `json:"type"`
	Text    string       `json:"text"`
	Payload interface{}  `json:"payload"`
}

type GameStart struct {
	GameId      uuid.UUID   `json:"game_id"`
	Color       Color       `json:"color"`
	TimeControl TimeControl `json:"time_control"`
}

type MovePiece struct {
	From   string `json:"from"`
	To     string `json:"to"`
	GameId string `json:"game_id" mapstructure:"game_id"`
}

type MoveResponse struct {
	From   string    `json:"from"`
	To     string    `json:"to"`
	Time   int64     `json:"time"`
	GameId uuid.UUID `json:"game_id"`
}

type GameOverResponse struct {
	Reason string    `json:"reason"`
	GameId uuid.UUID `json:"game_id"`
	Winner bool      `json:"winner"`
}

type TimeControl struct {
	Duration  string `json:"duration"`
	Increment string `json:"increment"`
}

type MatchParams struct {
	Players     []*Player   `json:"players"`
	TimeControl TimeControl `json:"time_control"`
}
