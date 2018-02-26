package Event

import (
	"github.com/GoGhost/pb"
	"github.com/GoGhost/rpc/client"
	"golang.org/x/net/context"
	"time"
)

type UserEvent interface {
	UserAddIn() error
	UserGetOut() error
	UserGrab() (string, error)
	UserGrabDone() error
	UserPlay() error
	UserPlayDone() error
}

type UserEventImpl struct {

}

var UE UserEvent = &UserEventImpl{}

func (impl *UserEventImpl) UserAddIn() error  {
	return nil
}

func (impl *UserEventImpl) UserGetOut() error  {
	return nil
}

func (impl *UserEventImpl) UserGrab() (string, error)  {

	cli, err := client.GRpcClientManager.Get("lord")

	if err != nil {
		return "", err
	}

	c, err := cli.Get("localhost:50051")
	if err != nil {
		return "", err
	}

	conn := pb.NewAskActionClient(c)

	ctx,cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var uid int64 = 1
	actionId := pb.ActionID_ActionGrab
	answer, err := conn.Ask(ctx, &pb.ActionRequest{Uid:&uid, Action:&actionId})
	if err != nil {
		return "", err
	}

	return *answer.Message, nil
}

func (impl *UserEventImpl) UserGrabDone() error  {
	return nil
}

func (impl *UserEventImpl) UserPlay() error  {
	return nil
}

func (impl *UserEventImpl) UserPlayDone() error  {
	return nil
}