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

import (
	"net"
	"time"
	log "github.com/GoGhost/log"
	"errors"
	"github.com/GoGhost/net/netInterface"
	"github.com/GoGhost/util"
	"github.com/GoGhost/protoCodec"
)

type TcpClient struct {

	remoteName string

	conn netInterface.Connection

	autoReConnect bool

	chanSize  uint32

	attributes *util.SyncMap

	netWorkCB netInterface.ClientCallBack

	remoteAddress string

	protocolName string

	timeInterval time.Duration
}

func NewTcpClient(name string, cs uint32, networkcb netInterface.ClientCallBack, protoName string, address string) netInterface.Client {
	client := &TcpClient{
		remoteName: name,
		chanSize: cs,
		autoReConnect: true,
		attributes: util.NewSyncMap(),
		netWorkCB: networkcb,
		remoteAddress: address,
		protocolName: protoName,
		timeInterval: 30 * time.Second,
	}

	log.Info("NewTcpClient complete ", client)
	return client
}

func (c *TcpClient) GetRemoteName() string {
	return c.remoteName
}

func (c *TcpClient) GetRemoteAddress() string {
	return c.remoteAddress
}

func (c *TcpClient) Start() bool {

	log.Info("Start Connect to ", c.remoteName, " address ", c.remoteAddress)
	tcpConn, err := net.DialTimeout("tcp", c.remoteAddress, 5*time.Second)
	if err != nil {
		log.Error(err.Error())
		c.Reconnect()
		return false
	}

	protocol := proto.NewProto(c.protocolName, tcpConn.(*net.TCPConn))

	c.conn = NewClientConn(GenConnId(), c.chanSize, c, protocol)

	//put into clientGroup
	c.conn.Start()

	return true
}

func (c *TcpClient) Stop() {

	c.autoReConnect = false
	
	if c.conn != nil {
		c.conn.Stop()
	}
}

func (c *TcpClient) GetConnection() netInterface.Connection {
	return c.conn
}

func (c *TcpClient)  SendMessage(msg netInterface.Message) error {
	if c.conn != nil && c.conn.Writable() {
		return c.conn.Write(msg)
	}

	return errors.New("TcpClient is not running")
}

func (c *TcpClient) SendData(data []byte) error {
	if c.conn != nil && c.conn.Writable() {
		_, err := c.conn.WriteBinary(data)
		return err
	}

	return errors.New("TcpClient is not running")
}


func (c *TcpClient) AutoReconnect() bool {
	return c.autoReConnect
}

func (this *TcpClient) SetAttribute(key string, attr interface{}) error {
	this.attributes.Put(key, attr)
	return nil
}

func (this *TcpClient) GetAttribute(key string) (attr interface{}, err error) {
	attr, ok := this.attributes.Get(key)

	if !ok {
		return nil, errors.New(key + " not exist")
	}

	return attr, nil
}

func (this *TcpClient) AllAttributes() (attrs map[interface{}]interface{}, err error) {
	return this.attributes.All(), nil
}

func (this *TcpClient) RemoveAttributes(key string) error {
	this.attributes.Remove(key)
	return nil
}

func (this *TcpClient) Reconnect() {
	time.AfterFunc(5*time.Second, func() {
		log.Error("Auto Reconnect server ", this.remoteName, " address ", this.remoteAddress)
		this.Start()
	})
}

//callback
func (this *TcpClient) OnMessageData(conn netInterface.Connection, msg netInterface.Message) error {

	if this.netWorkCB != nil {
		this.netWorkCB.OnMessageData(this, msg)
	}

	return nil
}

func (this *TcpClient) OnConnection(conn netInterface.Connection) {

	if this.netWorkCB != nil {
		this.netWorkCB.OnConnection(this)
	}

}

func (this *TcpClient) OnDisConnection(conn netInterface.Connection) {

	if this.netWorkCB != nil {
		this.netWorkCB.OnDisConnection(this)
	}
}

func (this *TcpClient) GetTimer() time.Duration {
	return this.timeInterval
}

func (this *TcpClient) StartTimer() {

	if this.conn != nil && this.conn.Writable() {
		//
		if this.netWorkCB != nil {
			this.netWorkCB.OnTimer(this)
		}
		//next
		time.AfterFunc(this.timeInterval, func() {
			this.StartTimer()
		})
	}
}