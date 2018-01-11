package websocket

import (
	"github.com/GoGhost/protoCodec"
	"github.com/GoGhost/net/netInterface"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

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

func NewWSProtocol(conn interface{}) netInterface.Codec {

	p := &wsProtocol{
		wsConn: conn.(*websocket.Conn),
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

func (c *wsProtocol) Read() (msg netInterface.Message, err error) {
	jmsg := &WsMsg{}

	buf := []byte{}
	var e error
	for {

		_, buf, e = c.wsConn.ReadMessage()

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

	return 0, c.wsConn.WriteMessage(2, pdubuf)
}

func (c *wsProtocol) WriteBinary(b []byte) (n int, err error) {
	return 0, c.wsConn.WriteMessage(2, b)
}

func (c *wsProtocol) Close() error {
	return c.wsConn.Close()
}

