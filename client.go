package main

import (
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

func (c *Client) Write() {
	for {
		msg, ok := <-c.Outgoing

		if ok {
			c.socket.WriteJSON(msg)
		}
	}
}

func (c *Client) Read() {
	for {
		defer c.socket.Close()

		var response Response
		c.socket.ReadJSON(&response)

		c.Incoming <- response
	}
}
