package main

import (
	"fmt"
	"time"
	"github.com/GoGhost/websocket"
 	_ "net/http/pprof"
	"net/http"
	"github.com/GoGhost/echo"
	"github.com/GoGhost/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pbh "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
	log "github.com/GoGhost/log"
	"net"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pbh.HelloRequest) (*pbh.HelloReply, error) {
	fmt.Println("Server recv : ", in.String())
	return &pbh.HelloReply{Message: "Hello " + in.Name}, nil
}

type LordActionService struct{}

func (las * LordActionService) waitPrepare(ctx context.Context, in *pb.ActionRequest) (*pb.ActionReply, error) {
	r := "ready"
	return &pb.ActionReply{Message: &r}, nil
}

func (las * LordActionService) waitGrab(ctx context.Context, in *pb.ActionRequest) (*pb.ActionReply, error) {
	time.Sleep(time.Second * 31)
	r := "ready"
	return &pb.ActionReply{Message: &r}, nil
}

func (las * LordActionService) waitPlay(ctx context.Context, in *pb.ActionRequest) (*pb.ActionReply, error) {
	r := "ready"
	return &pb.ActionReply{Message: &r}, nil
}


func (las * LordActionService) Ask(ctx context.Context, in *pb.ActionRequest) (*pb.ActionReply, error) {

	fmt.Println(in)

	switch in.GetAction() {
	case pb.ActionID_ActionPrepare:
		return las.waitPrepare(ctx, in)
	case pb.ActionID_ActionGrab:
		return las.waitGrab(ctx, in)
	case pb.ActionID_ActionPlay:
		return las.waitPlay(ctx, in)
	}

	return &pb.ActionReply{}, fmt.Errorf("no handler")
}

func main() {

	go func() {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterAskActionServer(s, &LordActionService{})

		// Register reflection service on gRPC server.
		reflection.Register(s)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	fmt.Println("start ECHO")

	echo.StartEchoServer()

	fmt.Println("start WS")

	websocket.StartWSServer()


	time.Sleep(1000 * time.Second)
}
