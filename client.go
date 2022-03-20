package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// Test Client
type Client struct {
	socket *websocket.Conn

	Incoming chan Response
	Outgoing chan Message
}

// Create a new test client
func NewClient() (*Client, error) {
	socket, _, err := websocket.DefaultDialer.Dial("ws://0.0.0.0:8080", nil)

	if err != nil {
		return nil, err
	}

	client := &Client{
		socket: socket,

		Incoming: make(chan Response),
		Outgoing: make(chan Message),
	}

	go client.Write()
	go client.Read()

	return client, nil
}

// Send data to the server
func (c *Client) Send(messageType MessageType) {
	c.Outgoing <- Message{
		Type: messageType,
	}
}

func (c *Client) Close() {
	c.socket.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(time.Second),
	)
}

func (c *Client) Write() {
	for {
		msg := <-c.Outgoing
		c.socket.WriteJSON(msg)
	}
}

func (c *Client) Read() {
	defer c.Close()

	for {
		var response Response
		err := c.socket.ReadJSON(&response)

		if err != nil {
			break
		}

		c.Incoming <- response
	}
}
