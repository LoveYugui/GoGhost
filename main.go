package main

import (
	_ "fmt"
	"github.com/GoGhost/echo"
	"time"
	"github.com/GoGhost/websocket"
)

func main() {
	echo.StartEchoServer()

	websocket.StartWebSocketServer()

	time.Sleep(1000 * time.Second)
}
