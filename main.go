package main

import (
	"fmt"
	"time"
	"github.com/GoGhost/websocket"
 	_ "net/http/pprof"
	"net/http"
	"github.com/GoGhost/echo"
)

func main() {

	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	fmt.Println("start ECHO")

	echo.StartEchoServer()

	fmt.Println("start WS")

	websocket.StartWSServer()


	time.Sleep(1000 * time.Second)
}
