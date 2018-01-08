package websocket

import (
	"github.com/GoGhost/net/netImplement/tcp"
	"github.com/GoGhost/net/netInterface"
	log "github.com/GoGhost/log"
	"net/http"
	"github.com/GoGhost/protoCodec"
)

type WebSocketServer struct {
	srv tcp.TCPServer
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

func StartLoop(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Client connection: %s", conn.RemoteAddr().String())

	protocol := proto.NewProto("WSProtocol", conn)

	if protocol == nil {
		return
	}

	tcpConnection:= tcp.NewServerConn(
		tcp.GenConnId(),
		conn,
		nil, //server,
		protocol)

	tcpConnection.Start()
	}

func StartWebSocketServer(w http.ResponseWriter, r *http.Request) {

	http.HandleFunc("/ws", StartLoop)

	if err := http.ListenAndServe("127.0.0.1:6667", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}