package echo

import (
	log "github.com/GoGhost/log"
	"github.com/GoGhost/net/netImplement/tcp"
	"github.com/GoGhost/net/netInterface"
)

type echoHandler struct {

}

func NewEchoHandler() netInterface.ConnectionCallBack {
	return &echoHandler{}
}

func (eh * echoHandler) OnConnection(conn netInterface.Connection) {
	log.Infoln("echoHandler OnConnection : ", conn)
}

func (eh * echoHandler) OnDisConnection(conn netInterface.Connection) {
	log.Infoln("echoHandler OnDisConnection : ", conn)
}

func (eh * echoHandler) OnMessageData(conn netInterface.Connection, msg netInterface.Message) error {
	log.Infoln("echoHandler OnMessageData : ", conn)
	conn.Write(msg)
	return nil
}

func StartEchoServer() {
	echoSrv, err := tcp.NewTCPServer(uint16(1), "127.0.0.1:6666", "EchoProtocol", NewEchoHandler())
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	echoSrv.Start()
}