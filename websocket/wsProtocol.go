package websocket

import (
	"github.com/GoGhost/protoCodec"
	"net"
	"github.com/GoGhost/net/netInterface"
	"fmt"
	"gitlab.mogujie.org/mgc/baselib/monitor"
	"log"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	proto.RegisterProto("WSProtocol", NewWSProtocol)
}

type WsMsg struct {
	buf []byte
}

func (msg *WsMsg) GetBuffer() ([]byte, error) {
	return msg.buf, nil
}


func (msg *WsMsg) HashCode() uint64 {
	return 0
}

func (msg *WsMsg) Encode() ([]byte, error) {
	return msg.buf, nil
}

type wsProtocol struct {
	wsConn *websocket.Conn
}

func NewWSProtocol(conn net.Conn) netInterface.Protocol {

	conn, err := upgrader.Upgrade(w, r, nil)

	var nc net.Conn = (net.Conn)(conn)

	p := &wsProtocol{
		wsConn: nc
	}

	return p
}

func (c *wsProtocol) GetType() string {
	return "EchoProtocol"
}

type echoMessage struct {
	data string
}

func NewEchoMessage() netInterface.Message {
	return &echoMessage{}
}

func (msg * echoMessage) GetBuffer() ([]byte, error) {
	return []byte{}, nil
}

func (msg * echoMessage) HashCode() uint64 {
	return 0
}

func (msg * echoMessage) Encode() ([]byte, error) {
	return ([]byte)(msg.data), nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *wsProtocol) Read(conn net.Conn) (msg netInterface.Message, err error) {
	jmsg := &WsMsg{}

	buf := []byte{}
	var e error
	for {

		var c *websocket.Conn =  conn.(*websocket.Conn)

		_, buf, e = c.(*websocket.Conn).ReadMessage()

		if e != nil {
			return nil, e
		}
		break

	}

	jmsg.buf = buf

	return jmsg, nil
}

func (c *wsProtocol) Write(msg netInterface.Message) (n int, err error) {

	pdubuf, err := msg.Encode()
	if err != nil {
		return 0, err
	}

	return c.TcpConn.Write(pdubuf)
}

func (c *wsProtocol) WriteBinary(b []byte) (n int, err error) {
	return c.TcpConn.Write(b)
}

func (c *wsProtocol) Close() error {
	return c.TcpConn.Close()
}

