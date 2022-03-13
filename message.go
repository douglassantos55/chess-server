package main

type MessageType string
type ResponseType string

const (
	QueueUp MessageType = "queue_up"
)

const (
	WaitForMatch ResponseType = "wait_for_match"
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
