package Game

import "github.com/GoGhost/MyLord/Role"
import (
	_ "github.com/GoGhost/log"
	"fmt"
	"log"
	_ "os"

	_ "golang.org/x/net/context"
	_ "google.golang.org/grpc"
	_ "google.golang.org/grpc/examples/helloworld/helloworld"
	"time"
	"github.com/GoGhost/rpc/client"
)

//
// protoc --go_out=plugins=grpc:. *.proto
//

const (
	PlayTimout = 30
)

const (
	GameStateUnstart = iota
	GameStateShuffle
	GameStateDeal
	GameStateGrab
	GameStatePlaying
	GameStateEnd
)

type LordGame struct {
	Roles [3]*Role.UserRole
	Cards [54]Role.CardRole
	State uint8
	Waiter *Waiter
	Notify chan uint8
}

func NewLordGame() *LordGame  {
	game := &LordGame{
		State:GameStateUnstart,
		Waiter:nil,
		Notify:make(chan uint8),
	}

	for i := 0; i < 3; i++ {
		game.Roles[i] = Role.NewUserRole()
	}

	for i := Role.SuitDiomand; i < Role.SuitSmallKing; i++ {
		for j := 0; j < 13; j++ {
			game.Cards[i*13 + j].Number = (uint8(j+Role.Number3))
			game.Cards[i*13 + j].Suit = (uint8(i))
			//fmt.Println(i*13 + j, game.Cards[i*13 + j].Number, game.Cards[i*13 + j].Suit)
		}
	}

	game.Cards[52].Number = Role.NumberSmallKing
	game.Cards[52].Suit = Role.SuitSmallKing

	game.Cards[53].Number = Role.NumberBigKing
	game.Cards[53].Suit = Role.SuitBigKing

	game.Waiter = NewWaiter(game)

	return game
}

func (game *LordGame) Start() error {
	if game == nil {
		return fmt.Errorf("game is nil")
	}

	if game.State != GameStateUnstart {
		return fmt.Errorf("game is already started")
	}

	game.State = GameStateShuffle

	err := game.Waiter.Shuffle()
	if err != nil {
		return err
	}

	//发牌
	game.State = GameStateDeal
	go game.Waiter.Deal()

	state := <- game.Notify
	if state != GameStateGrab {
		return fmt.Errorf("game state should be Grab here")
	}
	game.State = GameStateGrab

	for i := 0; i < 3; i++ {
		for j := 0; j < 17; j++ {
			fmt.Print(game.Roles[i].Cards[j].Number, "-",game.Roles[i].Cards[j].Suit," ")
		}
		fmt.Println()
	}


	//抢lord
	go game.Waiter.Grab()
	state = <- game.Notify
	if state != GameStatePlaying {
		return fmt.Errorf("game state should be Playing here")
	}
	game.State = GameStatePlaying

	/*
	//正式开始
	go game.Waiter.Play()
	state := <- game.Notify
	if state != GameStateEnd {
		return fmt.Errorf("game state should be End here")
	}
	game.State = GameStateEnd

	//计分统计
	*/
	return nil
}

//---------

const (
	address     = "localhost:50051"
	defaultName = "lord"
)

func Close()  {
	time.AfterFunc(time.Second * 20, func() {
		client.GRpcClientManager.Close(defaultName)
	})
}

func SayHello() {

	cli := client.NewRpcClient(defaultName)
	cli.Add(address)
	client.GRpcClientManager.RegisteRpcClient(defaultName, cli)

	cli, err := client.GRpcClientManager.Get(defaultName)

	if err != nil {
		log.Fatal("err : %s", err)
	}
	// Contact the server and print out its response.
	//name := defaultName
	//if len(os.Args) > 1 {
	//	name = os.Args[1]
	//}

	_, err = cli.Get(address)
	if err != nil {
		log.Fatal("err : %s", err)
	}

	for i := 0; i < 10000 ; i++ {
		//log.Printf("Greeting: state %s", conn.GetState().String())
		time.Sleep(time.Second * 2)
		/*
		conn.GetState()
		_, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		//log.Printf("Greeting: %s", r.Message)
		log.Printf("Greeting: end %s", time.Now().String())
		*/
	}
}

func init() {
	cli := client.NewRpcClient(defaultName)
	cli.Add(address)
	client.GRpcClientManager.RegisteRpcClient(defaultName, cli)

	cli, err := client.GRpcClientManager.Get(defaultName)

	if err != nil {
		log.Fatal("err : %s", err)
	}

	_, err = cli.Get(address)
	if err != nil {
		log.Fatal("err : %s", err)
	}
}