package proto

import (
	"net"
	"github.com/GoGhost/net/netInterface"
	"sync"
)

var protoFactory = make(map[string]func(conn net.Conn)netInterface.Protocol, 16)

var lock sync.Mutex

func RegisterProto(protoName string, f func(conn net.Conn) netInterface.Protocol) {
	lock.Lock()
	defer lock.Unlock()
	protoFactory[protoName] = f
}

func NewProto(protoName string, conn net.Conn) netInterface.Protocol {

	lock.Lock()
	defer lock.Unlock()

	if f, ok := protoFactory[protoName]; ok {
		return f(conn)
	}

	return nil
}
