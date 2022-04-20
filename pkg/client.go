package pkg

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

	close(c.Incoming)
	close(c.Outgoing)
}

func (c *Client) Write() {
	for {
		msg, ok := <-c.Outgoing

		if !ok {
			break
		}

		c.socket.WriteJSON(msg)
	}
}

func (c *Client) Read() {
	for {
		var response Response
		err := c.socket.ReadJSON(&response)

		if err != nil {
			break
		}

		c.Incoming <- response
	}
}
