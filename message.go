package main

type MessageType string
type ResponseType string

const (
	QueueUp      MessageType = "queue_up"
	Dequeue                  = "dequeue"
	Disconnected             = "disconnected"
)

const (
	WaitForMatch ResponseType = "wait_for_match"
	MatchFound                = "match_found"
)

type Message struct {
	Type   MessageType
	Text   string
	Player *Player
}

type Response struct {
	Type ResponseType
	Text string
}
