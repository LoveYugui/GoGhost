/*
 *  Copyright (c) 2018, https://github.com/LoveYugui/GoGhost
 *  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tcp

// #include <pthread.h>
// #include <unistd.h>
// #include <sys/syscall.h>
//
/*
int getpid() {
	return getpid();
}

long getself() {
	return pthread_self();
}
*/
import "C"

import (
	"flag"
	"sync"
	"errors"
	"sync/atomic"
	"fmt"
	"time"
	"runtime"
	"net"
	"github.com/GoGhost/net/netInterface"
	log "github.com/GoGhost/log"
	"github.com/GoGhost/constants"
	"github.com/GoGhost/global"
	"github.com/GoGhost/protoCodec"
)

func init() {
	flag.Parse()
	netIdentifier = 0
}

var (
	netIdentifier uint64
)

type TcpServerConfig struct {
	Name 		   string
	NetWorkVersion string
	Address        string // [ip:port] format
	MaxConnCount   uint32 // 最大连接数
	ProtocolName   string
}

func GetNetId() uint64 {
	return atomic.AddUint64(&netIdentifier, 1);
}

type TCPServer struct {

	srvNumber uint16

	isRunning bool

	address   string

	listener  net.Listener

	once      sync.Once

	connManager netInterface.ConnectionManager

	ProtocolName  string        // Protocol -> Make Codec

	NetworkCB netInterface.ConnectionCallBack // TcpConnection callBack

	serverConfig TcpServerConfig
}

func NewTCPServer(srvNum uint16, address string, protocolName string, netcb netInterface.ConnectionCallBack) (netInterface.Server, error) {

	server := &TCPServer{
		srvNumber: srvNum,
		isRunning: true,
		connManager: NewTcpConnectionManager(),
		once: sync.Once{},
		ProtocolName: protocolName,
		address: address,
		NetworkCB: netcb,
	}

	return server, nil
}

func (this *TCPServer)OnMessageData(conn netInterface.Connection, msg netInterface.Message) error {


	if this.NetworkCB != nil {
		return this.NetworkCB.OnMessageData(conn, msg)
	}

	log.Info("recv from " , conn, " msg : %v", msg)
	return nil
}

func (this *TCPServer)OnConnection(conn netInterface.Connection) {

	log.Infoln("------------------> connection add")

	if this.connManager != nil {
		this.connManager.AddConn(conn)
	}

	if this.NetworkCB != nil {
		this.NetworkCB.OnConnection(conn)
	}
}

func (this *TCPServer)OnDisConnection(conn netInterface.Connection) {

	if this.connManager != nil {
		this.connManager.RemoveConn(conn.GetId())
	}

	if this.NetworkCB != nil {
		this.NetworkCB.OnDisConnection(conn)
	}
}

func (server *TCPServer) Start() bool {

	exitChan := make(chan bool)

	go func(server *TCPServer) {

		defer func() {
			exitChan <- true
		}()

		runtime.LockOSThread()

		fmt.Println("StartAndServe inside")

		tcpAddr, err := net.ResolveTCPAddr("tcp4", server.address)
		if err != nil{
			log.Error("BAD CONFIG ", err, " content ", server.serverConfig)
			return
		}
		server.listener, err = net.ListenTCP("tcp4", tcpAddr)
		if err != nil {
			log.Error(err)
			return
		}

		log.Info("listen on ", *tcpAddr)

		for {
			var self int64 = int64(C.getself())

			log.Info(" pthread_self : ", self)

			conn, err := server.listener.Accept()
			if err != nil {
				log.Error(err)
				break
			}

			log.Info("new conn : ", conn)

			if server.connManager.Size() >= constants.MAX_CONNECTIONS {
				conn.Close()
			} else {
				server.establishTcpConnection(conn.(*net.TCPConn))
			}
		}

		log.Info("quit listen")
	}(server)


	<- exitChan

	log.Info("server stop ....")

	return true
}

func (this *TCPServer) Stop() bool {
	this.listener.Close()
	this.connManager.Dispose()

	return true
}

func GenConnId() uint64 {
	high := uint64(global.GetServerNumber())
	low  := GetNetId() & 0xffffffff
	return (high << 32) | low

}

func (server *TCPServer) establishTcpConnection(conn *net.TCPConn) {

	conn.SetReadDeadline(time.Now().Add(time.Minute*4))

	protocol := proto.NewProto(server.ProtocolName, conn)

	if protocol == nil {
		return
	}

	tcpConnection:= NewServerConn(
		GenConnId(),
		conn,
		server,
		protocol)

	tcpConnection.Start()
}

func (this *TCPServer) GetConfig() interface{} {
	return nil
}

func (this *TCPServer) GetManager() netInterface.ConnectionManager {
	return this.connManager
}

func (this *TCPServer) SendData(connId uint64, msg netInterface.Message) error {
	session := this.connManager.GetConn(connId)
	if session == nil {
		return errors.New("can not get session!")
	}

	return session.Write(msg)
}
