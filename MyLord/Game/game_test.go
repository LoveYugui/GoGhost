package Game

import (
	"testing"
	"fmt"
)

func TestNewLordGame(t *testing.T) {
	g := NewLordGame()

	//t.Log(g.Cards)

	fmt.Println("====")

	g.Start()
}

func TestSayHello(t *testing.T) {
	Close()
	SayHello()
}
