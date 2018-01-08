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
	log "github.com/GoGhost/log"
	"errors"
	"math/rand"
	"github.com/GoGhost/net/netInterface"
)

type TcpClientGroup struct {

	protocolName string

	clientMapLock sync.RWMutex
	clientMap map[string]map[string]netInterface.Client

	networkCallBack netInterface.ClientCallBack
}

func NewTcpClientGroup(protoName string, clients map[string][]string, cb netInterface.ClientCallBack) netInterface.ClientGroup  {
	group := &TcpClientGroup{
		protocolName: protoName,
		clientMap: make(map[string]map[string]netInterface.Client),
		networkCallBack: cb,
	}

	for k, v := range clients {

		m := make(map[string]netInterface.Client)

		for _, address := range v {
			client := NewTcpClient(k, 10 * 1024, group.networkCallBack, group.protocolName, address)

			if client != nil {
				m[address] = client
			}
		}

		group.clientMapLock.Lock()
		group.clientMap[k] = m
		group.clientMapLock.Unlock()

	}

	log.Info("NewTcpClientGroup group : ", group.clientMap)
	return group
}

func (group *TcpClientGroup) Start() bool {

	group.clientMapLock.Lock()
	defer group.clientMapLock.Unlock()

	for _, v := range group.clientMap {

		for _, c := range v {
			c.Start()
		}

	}

	return true
}

func (group *TcpClientGroup) Stop() bool {

	group.clientMapLock.Lock()
	defer group.clientMapLock.Unlock()

	for _, v := range group.clientMap {

		for _, c := range v {
			c.Stop()
		}

	}
	return true
}

func (group *TcpClientGroup) GetConfig() interface{} {
	return nil
}

func (group *TcpClientGroup) AddClient(name string, address string) {

	log.Info("TcpClientGroup AddClient name ", name, " address ", address)
	group.clientMapLock.Lock()
	defer group.clientMapLock.Unlock()

	m, ok := group.clientMap[name]

	if !ok {
		group.clientMap[name] = make(map[string]netInterface.Client)
	}

	m, _ = group.clientMap[name]

	_, ok = m[address]

	if ok {
		return
	}

	client := NewTcpClient(name, 10 * 1024, group.networkCallBack, group.protocolName, address)

	m[address] = client

	client.Start()
}

func (group *TcpClientGroup) RemoveClient(name string, address string) {

	log.Info("TcpClientGroup RemoveClient name ", name, " address ", address)

	group.clientMapLock.Lock()
	defer group.clientMapLock.Unlock()

	m, ok := group.clientMap[name]

	if !ok {
		return
	}

	m, _ = group.clientMap[name]

	c, ok := m[address]

	if !ok {
		return
	}

	c.Stop()

	delete(group.clientMap[name], address)

}

func (group *TcpClientGroup) SendData(name string, msg netInterface.Message) error {

	tcpConn := group.getRotationSession(name)

	if tcpConn == nil {
		return errors.New("Can not get connection!!")
	}

	//log.Info("yugui SendData ", name, tcpConn.RAddr().String(), " data :" , msg)

	return tcpConn.Write(msg)
}

func (this *TcpClientGroup) getRotationSession(name string) netInterface.Connection {

	all_conns := this.getTcpClientsByName(name)

	if all_conns == nil || len(all_conns) == 0 {
		return nil
	}

	index := rand.Int() % len(all_conns)

	return all_conns[index]
}

func (this *TcpClientGroup) BroadcastData (name string, msg netInterface.Message) error {

	all_conns := this.getTcpClientsByName(name)

	if all_conns == nil || len(all_conns) == 0 {
		return nil
	}

	for _, conn := range all_conns {
		conn.Write(msg)
	}

	return nil
}

func (this *TcpClientGroup) getTcpClientsByName(name string) []netInterface.Connection {

	var all_conns []netInterface.Connection

	this.clientMapLock.RLock()

	serviceMap, ok := this.clientMap[name]

	if !ok {
		this.clientMapLock.RUnlock()
		return nil
	}

	for _, c := range serviceMap {
		if c != nil && c.GetConnection() != nil && c.GetConnection().Writable() {
			all_conns = append(all_conns, c.GetConnection())
		}
	}

	this.clientMapLock.RUnlock()

	return all_conns
}
