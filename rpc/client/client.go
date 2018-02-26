package client

import (
	"google.golang.org/grpc"
	"sync"
	"fmt"
	"golang.org/x/net/context"
	log "github.com/GoGhost/log"
	"google.golang.org/grpc/connectivity"
	"time"
)

var (
	GRpcClientManager *RpcClientManager = &RpcClientManager{
		clients:make(map[string]*RpcClient),
	}
)

type RpcClientManager struct {
	clients map[string]*RpcClient
	lock sync.Mutex
}

func (m * RpcClientManager)RegisteRpcClient(name string, client *RpcClient) error {
	if client == nil {
		return fmt.Errorf("RpcClient is nil")
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.clients[name]; !ok {
		m.clients[name] = client
	} else {
		return fmt.Errorf("%s already exist", name)
	}

	return nil
}

func (m * RpcClientManager)Get(name string) (*RpcClient, error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.clients[name]; !ok {
		return nil, fmt.Errorf("not find %s", name)
	}

	return m.clients[name], nil
}

func (m * RpcClientManager)UnRegisteRpcClient(name string) error {

	m.lock.Lock()
	defer m.lock.Unlock()

	if client, ok := m.clients[name]; ok {
		delete(m.clients, name)
		client.Close()
	} else {
		return fmt.Errorf("%s not exist", name)
	}

	return nil
}

func (m * RpcClientManager)Close(name string) error {

	m.lock.Lock()
	defer m.lock.Unlock()

	for name, client := range m.clients {
		client.Close()
		delete(m.clients, name)
	}

	return nil
}

type RpcClient struct {
	close bool
	name string
	conns map[string]*grpc.ClientConn
	liveConns map[string]bool
	stop map[string]chan struct{}
	lock sync.Mutex
}

func NewRpcClient(name string) *RpcClient {
	rpcClient := &RpcClient{
		name:name,
		conns:make(map[string]*grpc.ClientConn),
		liveConns:make(map[string]bool),
		stop:make(map[string]chan struct{}),
	}

	return rpcClient
}

func (c *RpcClient) Close() {

	c.lock.Lock()
	defer c.lock.Unlock()

	c.close = true

	for _, s := range c.stop{
		close(s)
	}

	for addr, conn := range c.conns{
		conn.Close()
		delete(c.conns, addr)
	}

	c.conns = nil
	c.liveConns = nil
	c.stop = nil
}

func (c *RpcClient) Add(addr string) error {

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.close {
		return fmt.Errorf("name : %s already close", c.name)
	}

	if _, ok := c.conns[addr]; ok {
		return fmt.Errorf("name : %s, client %s already exist", c.name, addr)
	}

	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		fmt.Errorf("error occur on connecting to %s : err : %v", addr, err)
	}

	c.conns[addr] = conn

	s := make(chan struct{})
	c.stop[addr] = s

	go watchClient(c, addr, conn, s)

	return nil
}

func (c *RpcClient) Remove(addr string) error {

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.close {
		return fmt.Errorf("name : %s already close", c.name)
	}

	if _, ok := c.conns[addr]; !ok {
		return fmt.Errorf("name : %s, client %s not exist", c.name, addr)
	}

	conn := c.conns[addr]
	conn.Close()
	s := c.stop[addr]
	close(s)

	delete(c.conns, addr)
	delete(c.stop, addr)

	return nil
}

func (c *RpcClient) Get(addr string) (*grpc.ClientConn, error) {

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.close {
		return nil, fmt.Errorf("name : %s already close", c.name)
	}

	if _, ok := c.conns[addr]; !ok {
		return nil, fmt.Errorf("name : %s, client %s not exist", c.name, addr)
	}

	return c.conns[addr], nil
}

func (c *RpcClient) Ready(addr string) error {

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.close {
		return fmt.Errorf("name : %s already close", c.name)
	}

	if _, ok := c.conns[addr]; !ok {
		return fmt.Errorf("name : %s, client %s not exist", c.name, addr)
	}

	c.liveConns[addr] = true

	log.Infoln("%s : %s ready !", c.name, addr)
	fmt.Println("%s : %s ready !", c.name, addr)
	return nil
}

func (c *RpcClient) Reset(addr string) error {

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.close {
		return fmt.Errorf("name : %s already close", c.name)
	}

	if _, ok := c.conns[addr]; !ok {
		return fmt.Errorf("name : %s, client %s not exist", c.name, addr)
	}

	delete(c.liveConns, addr)

	log.Infoln("%s : %s reset !", c.name, addr)
	fmt.Println("%s : %s reset !", c.name, addr)

	return nil
}

func watchClient(c *RpcClient, addr string, conn * grpc.ClientConn, stop chan struct{}) {

	//name := c.name

	state := connectivity.Idle

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)

	for {

		if c == nil || c.close == true {
			break
		}
		//select {
		//case <-stop:
		//	fmt.Println("-----", c.close)
		//	log.Infof("%s : %s watchClient exit", name, addr)
		//	fmt.Printf("%s : %s watchClient exit", name, addr)
		//	return
		//case <-time.After(time.Millisecond * 3000):
		//	fmt.Println("lalalal")
		//	break
		//}

		pre := state
		stateChanged := conn.WaitForStateChange(ctx, state)

		cancel()

		s := conn.GetState()

		if !stateChanged {
			continue
		}

		//fmt.Println("stateChanged : ", stateChanged, " pre : ", state.String(), " now :" , s.String())

		if s == connectivity.Ready {
			c.Ready(addr)
		} else if pre == connectivity.Ready {
			c.Reset(addr)
		}

		state = s

	}
}