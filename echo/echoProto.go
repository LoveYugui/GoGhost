package echo

import (
	"net"
	"github.com/GoGhost/net/netInterface"
	"github.com/GoGhost/protoCodec"
)

func init() {
	proto.RegisterProto("EchoProtocol", NewEchoProtocol)
}

type echoProtocol struct {
	TcpConn net.Conn
}

func NewEchoProtocol(conn interface{}) netInterface.Codec {
	p := &echoProtocol{
		TcpConn: conn.(*net.TCPConn),
	}

	return p
}

func (c *echoProtocol) GetType() string {
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

func (c *echoProtocol) Read() (msg netInterface.Message, err error) {
	msg = NewEchoMessage()
	var buf = make([]byte, 256)
	n, err := c.TcpConn.Read(buf)

	if err != nil {
		return nil, err
	}

	msg.(*echoMessage).data = string(buf[:n])

	return msg, nil
}

func (c *echoProtocol) Write(msg netInterface.Message) (n int, err error) {

	pdubuf, err := msg.Encode()
	if err != nil {
		return 0, err
	}

	return c.TcpConn.Write(pdubuf)
}

func (c *echoProtocol) WriteBinary(b []byte) (n int, err error) {
	return c.TcpConn.Write(b)
}

func (c *echoProtocol) Close() error {
	return c.TcpConn.Close()
}

