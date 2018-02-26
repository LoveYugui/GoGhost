package Game

import (
	"fmt"
	"math/rand"
	"time"
	"sort"
	"github.com/GoGhost/MyLord/Role"
	"github.com/GoGhost/MyLord/Event"
)

type Waiter struct {
	prepare chan struct{}
	game 	*LordGame
}

func NewWaiter(g *LordGame) *Waiter {
	if g == nil {
		return nil
	}

	waiter := &Waiter{
		prepare : make(chan struct{}),
		game:g,
	}

	return waiter
}

func FreeWaiter(w *Waiter)  {
	if w == nil {
		return
	}

	close(w.prepare)
	w.game = nil
}

func (w *Waiter) check(state uint8) error {

	if w == nil {
		return fmt.Errorf("waiter is nil")
	}

	if w.game == nil {
		return fmt.Errorf("game is nil")
	}

	if w.game.State != state {
		return fmt.Errorf("game state should be %d", state)
	}

	return nil
}

func (w *Waiter) Shuffle() error {

	err := w.check(GameStateShuffle)

	if err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())

	count := rand.Int() % 5 + 2
	//fmt.Println("count : ", count)

	for ; count > 0; count-- {
		for sz := 54; sz > 1; sz-- {
			r := rand.Int() % sz
			swapIndex := sz - 1
			w.game.Cards[swapIndex], w.game.Cards[r] = w.game.Cards[r], w.game.Cards[swapIndex]
		}
	}

	for i := 0; i < 54; i++ {
		//fmt.Print(w.game.Cards[i].Number, " ")
	}
	//fmt.Println("----")

	return nil
}

func (w *Waiter) Deal() error {

	err := w.check(GameStateDeal)

	if err != nil {
		return err
	}

	for i := 0; i < 17; i++ {
		for j := 0; j < 3; j++ {
			w.game.Roles[j].Cards[i] = w.game.Cards[i*3 + j]
		}
	}

	for i := 0; i < 3; i++ {
		sort.Sort(Role.CardSlice(w.game.Roles[i].Cards))    // 排序
	}

	w.game.Notify <- GameStateGrab
	return nil
}

func (w *Waiter) Grab() error {
	err := w.check(GameStateGrab)

	if err != nil {
		return err
	}

	luckyOne := rand.Int() % 3

	fmt.Println(luckyOne)

	for i := 0; i < 3; i++ {
		grabResp, err := Event.UE.UserGrab()
		if err != nil {
			fmt.Println("no grab ", err)
		} else {
			fmt.Println(grabResp)
		}
	}

	w.game.Notify <- GameStatePlaying
	return nil

}