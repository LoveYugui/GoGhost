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

package netInterface

import (
	"net"
	"time"
)

type Connection interface {
	Name() string
	GetId() uint64
	GetConn() net.Conn
	GetProtocol() Protocol
	Type() uint8
	State() uint8
	Start() bool
	Stop() bool
	Close()

	NeedReconnect() bool

	SetExtraData(d interface{})
	GetExtraData() interface{}

	SetTimeStamp(t time.Time)
	GetTimeStamp() time.Time

	//========
	Read() error
	Write(message Message) error
	Readable() bool
	Writable() bool
	WriteBinary(msg []byte) (n int, err error)

	SetAttribute(key string, attr interface{}) error
	GetAttribute(key string) (attr interface{}, err error)
	AllAttributes() (attrs map[interface{}]interface{}, err error)
	RemoveAttributes(key string) error
	LAddr() net.Addr
	RAddr() net.Addr

	//SetCallBack(actionName string, cb ConnectionCallBack)
}


type ConnectionCallBack interface {
	OnMessageData(conn Connection, msg Message) error
	OnConnection(conn Connection)
	OnDisConnection(conn Connection)
}

