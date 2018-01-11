
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
	"sync"
	"time"
	"fmt"
	log "github.com/GoGhost/log"
	"github.com/GoGhost/util"
	"net"
	"errors"
	"github.com/GoGhost/net/netInterface"
)

const (
	NTYPE = 4
	NLEN = 4
	MAXLEN = 1 << 23  // 8M
)

//socket state
const (
	CLOSED = iota
	CONNECTING
	ESTABLISHED
	LISTEN
	SOCKSTATE_NUM
)

//socket type
const (
	TCPLISTEN = iota
	TCPSVRCONN
	TCPCLIENTCONN
	UDPSOCK
	SOCKTYPE_NUM
)

type TcpConnection struct {

	ConnID             uint64
	name               string
	address            string
	connState          uint8
	connType           uint8

	//conn 		       *net.TCPConn

	codec	           netInterface.Codec

	running            *util.AtomicBoolean
	once               sync.Once
	finish             sync.WaitGroup

	heartBeat          bool
	heartBeatInterval  time.Duration

	//异步数据发送队列
	messageSendChan    chan netInterface.Message
	messageHandlerChan chan netInterface.Message
	closeConnChan      chan struct{}

	//业务事件的回调
	NetworkCB          netInterface.ConnectionCallBack
	Reconnect          bool

	//时间戳, 用作RT计算
	TimeStamp          time.Time

	//Work协程的数目
	WorkNum            int

	attributes 		   *util.SyncMap

	//每分钟接收数据包统计
	LastStatMinute	   int
	RecvCountInMinute  int32
}

type TimeLoopFun  func(conn *TcpConnection)

func NewServerConn(connid uint64, networkcb netInterface.ConnectionCallBack, p netInterface.Codec) netInterface.Connection {
	serverConn := &TcpConnection{
		ConnID: connid,

		running: util.NewAtomicBoolean(true),
		connState: CLOSED,
		connType: TCPSVRCONN,

		heartBeat : false,
		heartBeatInterval : 0,

		codec: p,
		finish: sync.WaitGroup{},
		messageSendChan: make(chan netInterface.Message, 128),
		messageHandlerChan: make(chan netInterface.Message, 128),
		closeConnChan: make(chan struct{}),

		NetworkCB: networkcb,
		Reconnect : false,

		WorkNum : 1,

		attributes: util.NewSyncMap(),
	}

	//serverConn.NetworkCB.OnConnection(serverConn)

	return serverConn
}

func NewClientConn(connId uint64, chanSize uint32, networkcb netInterface.ConnectionCallBack, p netInterface.Codec) (netInterface.Connection) {
	clientConn := &TcpConnection{
		ConnID: connId,

		running: util.NewAtomicBoolean(true),

		connState: CLOSED,
		connType: TCPCLIENTCONN,

		heartBeat : true,
		heartBeatInterval : 30 * time.Second,

		codec: p,
		finish: sync.WaitGroup{},
		messageSendChan: make(chan netInterface.Message, chanSize),
		messageHandlerChan: make(chan netInterface.Message, chanSize),
		closeConnChan: make(chan struct{}),
		NetworkCB: networkcb,
		WorkNum : 1,
		Reconnect : true,

		attributes: util.NewSyncMap(),
	}
	return clientConn
}

func (this *TcpConnection) String() string{
	return fmt.Sprintln(this.name, " , id = ", this.ConnID)//, " , lAddr : ", this.conn.LocalAddr(), " , rAddr : ", this.conn.RemoteAddr())
}

func (this *TcpConnection) SetTimeStamp(t time.Time) {
	this.TimeStamp = t
}

func (this *TcpConnection) GetTimeStamp() time.Time{
	return this.TimeStamp
}

func (this *TcpConnection) IsRunning() bool {
	return this.running.Get()
}

func (this *TcpConnection)SetWorkNum(num int) {
	this.WorkNum = num
}

func (this *TcpConnection)SetAddress(addr string) {
	this.address = addr
}

func (this *TcpConnection)Address() string {
	return this.address
}

func (this *TcpConnection) AddTimeTask(interval time.Duration, fn func()) {
}

func (this *TcpConnection)SetName(name string) {
	this.name = name
}

// implement of netInterface.IConn

func (this *TcpConnection)Name() string {
	return this.name
}

func (this *TcpConnection)GetId() uint64 {
	return this.ConnID
}

func (this *TcpConnection) GetConn() net.Conn  {
	return nil//this.conn
}

func (this *TcpConnection) GetCodec() netInterface.Codec {
	return this.codec
}

func (this *TcpConnection)Type() uint8 {
	return this.connType
}

func (this *TcpConnection)State() uint8 {
	return this.connState
}

func (this *TcpConnection)Start() bool {

	this.connState = ESTABLISHED

	this.finish.Add(2 + this.WorkNum)
	go this.readLoop()
	go this.writeLoop()
	for i := 0; i < this.WorkNum; i++ {
		go this.workLoop()
	}

	if this.NetworkCB != nil {
		this.NetworkCB.OnConnection(this)
	}
	return true
}

func (this *TcpConnection)Stop() bool {
	this.Close()
	return true
}

func (this *TcpConnection)Close() {
	this.once.Do(func() {
		if this.running.CompareAndSet(true, false) {

			this.connState = CLOSED

			//通知给上一层关闭
			if this.NetworkCB != nil {
				this.NetworkCB.OnDisConnection(this)
			}
			//关闭掉channel，所有的work都要监测closeConnChan
			close(this.messageSendChan)
			close(this.messageHandlerChan)
			close(this.closeConnChan)

			//关闭网络连接；
			this.codec.Close()
			this.finish.Wait()
		}
	})
}

func (this *TcpConnection)Readable() bool {
	return this.IsRunning()
}

func (this *TcpConnection)Writable() bool {
	return this.IsRunning()
}

func (this *TcpConnection)NeedReconnect() bool {
	return this.Reconnect
}

func (this *TcpConnection)Write(msg netInterface.Message) (err error) {

	if !this.IsRunning() {
		return errors.New("TcpConnection is not running")
	}

	select {
	case this.messageSendChan <- msg:
		return nil
	default:
		log.Error("messageSendChan is full , Write Lost packet !!! ", this.ConnID)
		return nil
	}
}

func (this *TcpConnection)WriteBinary(msg []byte) (n int, err error) {

	if !this.IsRunning() {
		return 0, errors.New("TcpConnection is not running")
	}

	return this.codec.WriteBinary(msg)
}

func (this *TcpConnection) Read()  error {

	msg, err := this.codec.Read()
	if err != nil {
		log.Warning("ReadData fail address")//, this.LAddr().String())
		return err
	}

	if msg == nil {
		//stats rtt
		//rtt := uint64(time.Since(this.GetTimeStamp()).Nanoseconds())  / uint64(time.Millisecond)
	} else {
		this.messageHandlerChan <- msg
	}

	return nil
}

func (this *TcpConnection) LAddr() net.Addr {
	return nil//this.conn.LocalAddr()
}

func (this *TcpConnection) RAddr() net.Addr {
	return nil//this.conn.RemoteAddr()
}

func (this *TcpConnection) SetAttribute(key string, attr interface{}) error {
	this.attributes.Put(key, attr)
	return nil
}

func (this *TcpConnection) GetAttribute(key string) (attr interface{}, err error) {
	attr, ok := this.attributes.Get(key)

	if !ok {
		return nil, errors.New(key + " not exist")
	}

	return attr, nil
}

func (this *TcpConnection) AllAttributes() (attrs map[interface{}]interface{}, err error) {
	return this.attributes.All(), nil
}

func (this *TcpConnection) RemoveAttributes(key string) error {
	this.attributes.Remove(key)
	return nil
}


func (this *TcpConnection)Init(name string, address string) bool {
	this.name = name
	this.address = address

	return true
}

func (this *TcpConnection)NetworkOwner() netInterface.ConnectionCallBack {
	return this.NetworkCB
}

func (this *TcpConnection)readLoop() {

	defer func () {
		//log.Error("readLoop error catched, close client")
		this.finish.Done()
		this.Close()
	}()

	for this.IsRunning() {
		select {
		case <-this.closeConnChan:
			log.Info("readLoop -> To Close ", this.String())
			return
		default:
			err := this.Read()
			if err != nil {
				log.Info("ReadLoop -> To Close:", err.Error() , " ", this.String())
				return
			}
		}
	}
}

func (this *TcpConnection)writeLoop() {

	defer func () {
		//log.Error("writeLoop error catched, close client")
		this.finish.Done()
		this.Close()
	}()

	for this.IsRunning() {
		select {
		case <-this.closeConnChan:
			log.Info("writeLoop -> To Close", this.String())
			return

		case msg := <-this.messageSendChan:
			if msg != nil {
				if _, err := this.codec.Write(msg); err != nil {
					log.Error("Error writing data ", err.Error(), " ", this.String())
					return
				}
			}
		}
	}
}

func (this *TcpConnection)workLoop() {

	defer func () {
		//log.Error("readLoop error catched, close client")
		this.finish.Done()
		this.Close()
	}()

	for this.IsRunning() {
		select {
		case <-this.closeConnChan:
			log.Info("workLoop -> To Close ", this.String())
			return

		case msg := <-this.messageHandlerChan:
			if msg != nil && this.NetworkCB != nil {
				if err := this.NetworkCB.OnMessageData(this, msg); err != nil {
				}
			}
		}
	}
}
