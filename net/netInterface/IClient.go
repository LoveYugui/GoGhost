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

import "time"

type Client interface {

	GetRemoteName() string

	GetRemoteAddress() string

	Start() bool

	Stop()

	GetConnection() Connection

	SendData(data []byte) error

	SendMessage(msg Message) error

	AutoReconnect() bool

	SetAttribute(key string, attr interface{}) error

	GetAttribute(key string) (attr interface{}, err error)

	AllAttributes() (attrs map[interface{}]interface{}, err error)

	RemoveAttributes(key string) error

	Reconnect()

	GetTimer() time.Duration

	StartTimer()
}

type ClientCallBack interface {
	OnMessageData(client Client, msg Message) error
	OnConnection(client Client)
	OnDisConnection(client Client)
	OnTimer(client Client)
}