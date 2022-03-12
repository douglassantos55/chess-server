package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	socket *websocket.Conn
}

func NewClient() (*Client, error) {
	socket, _, err := websocket.DefaultDialer.Dial("ws://0.0.0.0:8080", nil)

	if err != nil {
		return nil, err
	}

	return &Client{socket}, nil
}
