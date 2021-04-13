// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import (
	"errors"
	"fmt"
	"sync"
)

//              Server                                                            Client
//
// +---------------------------------+                               +---------------------------------+
// | SharedMemoryTransport           |                               | SharedMemoryTransport           |
// +----+-----------------------+----+                               +----+-----------------------+----+
//      |                       ^                                         |                       ^
//      |                       |                                         |                       |
//      |                       |                                         |                       |
//      v                       |                                         v                       |
// +----+-----------------------+----+                               +----+-----------------------+----+
// | SharedMemoryTransportConnection |                               | SharedMemoryTransportConnection |
// +----+-----------------------+----+                               +----+-----------------------+----+
//      |                       ^                                         |                       ^
//      |                       |                                         |                       |
//      |                       |                                         |                       |
//      |                       |         +---------------------+         |                       |
//      |                       |         | sharedMemoryChannel |         |                       |
//      |                       |         |                     |         |                       |
//      |                       +---------+  serverRxClientTx   +<--------+                       |
//      +-------------------------------->+  serverTxClientRx   +---------------------------------+
//                                        +---------------------+

// sharedMemoryChannel ...
type sharedMemoryChannel struct {
	serverTxClientRx chan *[]byte
	serverRxClientTx chan *[]byte
}

// newSharedMemoryChannel ...
func newSharedMemoryChannel(size int) *sharedMemoryChannel {
	return &sharedMemoryChannel{
		serverTxClientRx: make(chan *[]byte, size),
		serverRxClientTx: make(chan *[]byte, size),
	}
}

// SharedMemoryTransportConnection ...
type SharedMemoryTransportConnection struct {
	name      string
	forServer bool
	closed    bool
	channel   *sharedMemoryChannel
}

// NewSharedMemoryTransportConnection ...
func NewSharedMemoryTransportConnection(
	name string,
	forServer bool,
	channel *sharedMemoryChannel) TransportConnection {

	return &SharedMemoryTransportConnection{
		name:      name,
		forServer: forServer,
		closed:    false,
		channel:   channel,
	}
}

// GetName ...
func (c *SharedMemoryTransportConnection) GetName() string {
	return c.name
}

// Read ...
func (c *SharedMemoryTransportConnection) Read(p []byte) (n int, err error) {
	if c.closed {
		fmt.Println("connection is closed")
		return 0, nil
	}

	var data *[]byte

	if c.forServer {
		data = <-c.channel.serverRxClientTx
	} else {
		data = <-c.channel.serverTxClientRx
	}

	return copy(p, *data), nil
}

// Write ...
func (c *SharedMemoryTransportConnection) Write(p []byte) (n int, err error) {
	if c.closed {
		fmt.Println("connection is closed")
		return 0, nil
	}

	data := make([]byte, len(p))
	n = copy(data, p)

	if n != len(p) {
		return 0, errors.New("short write")
	}

	if c.forServer {
		c.channel.serverTxClientRx <- &data
	} else {
		c.channel.serverRxClientTx <- &data
	}

	return n, nil
}

// Close ...
func (c *SharedMemoryTransportConnection) Close() {
	c.closed = true
}

// Closed ...
func (c *SharedMemoryTransportConnection) Closed() bool {
	return c.closed
}

// SharedMemoryTransport ...
type SharedMemoryTransport struct {
	name       string
	address    string
	forServer  bool
	connection chan TransportConnection
}

// NewSharedMemoryTransport ...
func NewSharedMemoryTransport(name string) Transport {
	return &SharedMemoryTransport{
		name:       name,
		forServer:  true,
		connection: make(chan TransportConnection, 1),
	}
}

// ForServer ...
func (sm *SharedMemoryTransport) ForServer(forServer bool) Transport {
	sm.forServer = forServer
	return sm
}

// GetName ...
func (sm *SharedMemoryTransport) GetName() string {
	return sm.name
}

// SetAddress ...
func (sm *SharedMemoryTransport) SetAddress(address string) Transport {
	sm.address = address
	return sm
}

// GetAddress ...
func (sm *SharedMemoryTransport) GetAddress() string {
	return sm.address
}

// NextConnection ...
func (sm *SharedMemoryTransport) NextConnection() TransportConnection {
	return <-sm.connection
}

// Two Transports(Client side and Server side) should use the same SharedMemoryChannel.
// Keep a global list for reusing
var mutex sync.Mutex
var chans map[string]*sharedMemoryChannel = make(map[string]*sharedMemoryChannel, 32)

// Run ...
func (sm *SharedMemoryTransport) Run() Transport {
	mutex.Lock()
	defer mutex.Unlock()

	// check if some Transport have created
	ch, ok := chans[sm.address]
	if !ok {
		// not found, create new one
		ch = newSharedMemoryChannel(512)
		chans[sm.address] = ch
		fmt.Println(sm.address, ch)
	}
	sm.connection <- NewSharedMemoryTransportConnection(sm.name, sm.forServer, ch)

	return sm
}
