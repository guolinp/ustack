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
	refCount         int
	serverTxClientRx chan *[]byte
	serverRxClientTx chan *[]byte
}

// newSharedMemoryChannel ...
func newSharedMemoryChannel(size int) *sharedMemoryChannel {
	return &sharedMemoryChannel{
		refCount:         0,
		serverTxClientRx: make(chan *[]byte, size),
		serverRxClientTx: make(chan *[]byte, size),
	}
}

// deleteSharedMemoryChannel ...
func deleteSharedMemoryChannel(smc *sharedMemoryChannel) {
	if smc != nil {
		close(smc.serverTxClientRx)
		close(smc.serverRxClientTx)
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
	name    string
	options map[string]interface{}
	sync.Mutex
	address    string
	isRunning  bool
	forServer  bool
	connection TransportConnection
	next       chan TransportConnection
	queueSize  int
}

// NewSharedMemoryTransport ...
func NewSharedMemoryTransport(name string) Transport {
	return &SharedMemoryTransport{
		name:       name,
		options:    make(map[string]interface{}),
		isRunning:  false,
		forServer:  true,
		connection: nil,
		next:       make(chan TransportConnection, 1),
		queueSize:  512,
	}
}

// parseOptions ...
func (sm *SharedMemoryTransport) parseOptions() {
	size, exists := OptionParseInt(sm.GetOption("MaxQueueSize"), 512)
	sm.queueSize = size
	if exists {
		fmt.Println("SharedMemoryTransport: option MaxQueueSize:", sm.queueSize)
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

// SetOption ...
func (sm *SharedMemoryTransport) SetOption(name string, value interface{}) Transport {
	sm.options[name] = value
	return sm
}

// GetOption ...
func (sm *SharedMemoryTransport) GetOption(name string) interface{} {
	if value, ok := sm.options[name]; ok {
		return value
	}
	return nil
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
	next := <-sm.next
	sm.connection = next
	return sm.connection
}

// Two Transports(Client side and Server side) should use the same
// sharedMemoryChannel instance.
// Keep a global list for reusing
var mutex sync.Mutex
var chans map[string]*sharedMemoryChannel = make(map[string]*sharedMemoryChannel, 32)

// Run ...
func (sm *SharedMemoryTransport) Run() Transport {
	sm.Lock()
	defer sm.Unlock()

	if sm.isRunning {
		return sm
	}

	sm.parseOptions()

	mutex.Lock()

	// check if some Transport have created
	ch, ok := chans[sm.address]
	if !ok {
		// not found, create new one
		ch = newSharedMemoryChannel(sm.queueSize)
		chans[sm.address] = ch
	}

	ch.refCount++

	sm.next <- NewSharedMemoryTransportConnection(sm.name, sm.forServer, ch)

	mutex.Unlock()

	return sm
}

// Stop ...
func (sm *SharedMemoryTransport) Stop() Transport {
	sm.Lock()
	defer sm.Unlock()

	if !sm.isRunning {
		return sm
	}

	sm.connection.Close()

	mutex.Lock()

	ch, ok := chans[sm.address]
	if ok {
		ch.refCount--

		if ch.refCount == 0 {
			delete(chans, sm.address)
			deleteSharedMemoryChannel(ch)
		}
	}

	mutex.Unlock()

	return sm
}
