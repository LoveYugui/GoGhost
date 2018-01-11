package websocket

import (
	"github.com/GoGhost/net/netImplement/tcp"
	"github.com/GoGhost/net/netInterface"
	log "github.com/GoGhost/log"
	"net/http"
	"github.com/GoGhost/protoCodec"
	"sync"
	"fmt"
	"github.com/GoGhost/net/netImplement/manager"
)

type WebSocketServer struct {

	srvNumber uint16

	isRunning bool

	address   string

	once      sync.Once

	connManager netInterface.ConnectionManager

	ProtocolName  string        // Protocol -> Make Codec

	NetworkCB netInterface.ConnectionCallBack // TcpConnection callBack

}

type wsHandler struct {

}

func NewWSHandler() netInterface.ConnectionCallBack {
	return &wsHandler{}
}

func (eh * wsHandler) OnConnection(conn netInterface.Connection) {
	log.Infoln("wsHandler OnConnection : ", conn)
}

func (eh * wsHandler) OnDisConnection(conn netInterface.Connection) {
	log.Infoln("wsHandler OnDisConnection : ", conn)
}

func (eh * wsHandler) OnMessageData(conn netInterface.Connection, msg netInterface.Message) error {
	log.Infoln("wsHandler OnMessageData : ", conn)
	conn.Write(msg)
	return nil
}

func (server *WebSocketServer) startWSClient(w http.ResponseWriter, r *http.Request) {

	log.Infof("startWSClient")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("Client connection: %s", conn.RemoteAddr().String())

	protocol := proto.NewProto(server.ProtocolName, conn)

	if protocol == nil {
		return
	}

	tcpConnection:= tcp.NewServerConn(
		tcp.GenConnId(),
		server,
		protocol)

	tcpConnection.Start()
}

func NewWSServer(srvNum uint16, address string, protocolName string, m netInterface.ConnectionManager, netcb netInterface.ConnectionCallBack) (netInterface.Server, error) {

	server := &WebSocketServer{
		srvNumber: srvNum,
		isRunning: true,
		connManager: m,
		once: sync.Once{},
		ProtocolName: protocolName,
		address: address,
		NetworkCB: netcb,
	}

	return server, nil
}

func (this *WebSocketServer)OnMessageData(conn netInterface.Connection, msg netInterface.Message) error {


	if this.NetworkCB != nil {
		return this.NetworkCB.OnMessageData(conn, msg)
	}

	log.Info("recv from " , conn, " msg : %v", msg)

	return nil
}

func (this *WebSocketServer)OnConnection(conn netInterface.Connection) {

	log.Infoln("------------------> connection add")

	if this.connManager != nil {
		this.connManager.AddConn(conn)
	}

	if this.NetworkCB != nil {
		this.NetworkCB.OnConnection(conn)
	}
}

func (this *WebSocketServer)OnDisConnection(conn netInterface.Connection) {

	if this.connManager != nil {
		this.connManager.RemoveConn(conn.GetId())
	}

	if this.NetworkCB != nil {
		this.NetworkCB.OnDisConnection(conn)
	}
}

func (server *WebSocketServer) Start() bool {

	exitChan := make(chan bool)

	go func(server *WebSocketServer) {

		defer func() {
			exitChan <- true
		}()

		fmt.Println("StartAndServe inside")

		http.HandleFunc("/ws", server.startWSClient)

		if err := http.ListenAndServe(server.address, nil); err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}(server)


	<- exitChan

	log.Info(" ws server stop ....")

	return true
}

func (this *WebSocketServer) Stop() bool {
	return true
}

func (this *WebSocketServer) GetConfig() interface{} {
	return nil
}

func (this *WebSocketServer) GetManager() netInterface.ConnectionManager {
	return this.connManager
}

func StartWSServer()  {

	wsSrv, err := NewWSServer(uint16(2), "0.0.0.0:6667", "WSProtocol", manager.CM, NewWSHandler())
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	go wsSrv.Start()
}