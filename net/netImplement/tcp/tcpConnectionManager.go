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
	"math/rand"
	"time"
	"github.com/GoGhost/net/netInterface"
)

const sessionMapNum = 8

type TcpConnectionManager struct {
	sessionMaps [sessionMapNum]sessionMap
	disposeFlag bool
	disposeOnce sync.Once
	disposeWait sync.WaitGroup
	current     uint64
}

type sessionMap struct {
	sync.RWMutex
	sessions map[uint64]netInterface.Connection
}

func (this *sessionMap) RandomSelectFromMap() netInterface.Connection {
	var array_conns []netInterface.Connection
	for _, session := range this.sessions {
		array_conns = append(array_conns, session)
	}

	if len(array_conns) == 0 {
		return nil
	}

	index := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(array_conns))

	return array_conns[index]
}

func NewTcpConnectionManager() netInterface.ConnectionManager {
	manager := &TcpConnectionManager{}
	for i := 0; i < len(manager.sessionMaps); i++ {
		manager.sessionMaps[i].sessions = make(map[uint64]netInterface.Connection)
	}

	return manager
}

func (this *TcpConnectionManager) Dispose() {
	this.disposeOnce.Do(func() {
		this.disposeFlag = true
		for i := 0; i < sessionMapNum; i++ {
			smap := &this.sessionMaps[i]
			smap.Lock()
			for _, session := range smap.sessions {
				session.Close()
			}
			smap.Unlock()
		}
		this.disposeWait.Wait()
	})
}

func (this *TcpConnectionManager) GetConn(connId uint64) netInterface.Connection {
	smap := &this.sessionMaps[connId % sessionMapNum]
	smap.RLock()
	defer smap.RUnlock()
	session, _ := smap.sessions[connId]
	return session
}

// 查找所有远端为address客户端
func (this *TcpConnectionManager) GetConnsByAddress(address string)  (conns []netInterface.Connection) {
	this.BroadcastRun(func (conn netInterface.Connection) {
		if conn.RAddr().String() == address {
			conns = append(conns, conn)
		}
	})

	return conns
}

//通过轮训的方式找到hash
func (this *TcpConnectionManager) GetRotationConn() netInterface.Connection {
	var all_conns []netInterface.Connection

	this.BroadcastRun(func (conn netInterface.Connection){
		all_conns = append(all_conns, conn)
	})

	if len(all_conns) == 0 {
		return nil
	}

	index := (this.current) % uint64(len(all_conns))
	this.current++

	return all_conns[index]
}

func (this *TcpConnectionManager) GetRandomConn() netInterface.Connection {

	index := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(sessionMapNum)

	for i := 0; i < sessionMapNum; i++ {
		index++
		if index >= sessionMapNum {
			index = 0
		}

		smap := &this.sessionMaps[index]
		smap.RLock()
		conn := smap.RandomSelectFromMap()
		if conn != nil {
			smap.RUnlock()
			return conn
		}
		smap.RUnlock()
	}

	return nil
}

func (this *TcpConnectionManager) GetHashConn(hashCode uint64) netInterface.Connection {
	return nil
}

func (this *TcpConnectionManager) BroadcastRun(handler func(netInterface.Connection)) {

	for i := 0; i < sessionMapNum; i++ {
		smap := &this.sessionMaps[i]

		smap.RLock()
		for _, v := range smap.sessions {
			handler(v)
		}
		smap.RUnlock()
	}
}

func (this *TcpConnectionManager) AddConn(conn netInterface.Connection) {
	smap := &this.sessionMaps[conn.GetId() % sessionMapNum]
	smap.Lock()
	defer smap.Unlock()
	smap.sessions[conn.GetId()] = conn
	this.disposeWait.Add(1)
}

func (this *TcpConnectionManager) RemoveConn(connId uint64) {

	if this.disposeFlag {
		this.disposeWait.Done()
		return
	}
	smap := &this.sessionMaps[connId % sessionMapNum]
	smap.Lock()
	defer smap.Unlock()
	delete(smap.sessions, connId)
	this.disposeWait.Done()
}

func (this *TcpConnectionManager) AllConn() (conns []netInterface.Connection) {
	this.BroadcastRun(func (conn netInterface.Connection) {
			conns = append(conns, conn)
	})

	return conns
}

func (this *TcpConnectionManager) Size() uint32 {
	var s uint32 = 0
	for i := 0; i < sessionMapNum; i++ {
		smap := &this.sessionMaps[i]
		smap.RLock()
		s += uint32(len(smap.sessions))
		smap.RUnlock()
	}

	return s
}