package main

import (
	"time"

	"example.com/chess-server/pkg"
)

func main() {
	server := pkg.NewServer([]pkg.Handler{
		pkg.NewQueueManager(),
		pkg.NewMatchMaker(10 * time.Second),
		pkg.NewGameManager(),
	})
	server.Listen("0.0.0.0:8080")
}
