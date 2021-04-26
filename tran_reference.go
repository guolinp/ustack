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
// |   ReferenceTransport            |                               |   ReferenceTransport            |
// +----+-----------------------+----+                               +----+-----------------------+----+
//      |                       ^                                         |                       ^
//      |                       |                                         |                       |
//      |                       |                                         |                       |
//      v                       |                                         v                       |
// +----+-----------------------+----+                               +----+-----------------------+----+
// |   ReferenceTransportConnection  |                               |   ReferenceTransportConnection  |
// +----+-----------------------+----+                               +----+-----------------------+----+
//      |                       ^                                         |                       ^
//      |                       |                                         |                       |
//      |                       |                                         |                       |
//      |                       |         +---------------------+         |                       |
//      |                       |         |  referenceChannel   |         |                       |
//      |                       |         |                     |         |                       |
//      |                       +---------+  serverRxClientTx   +<--------+                       |
//      +-------------------------------->+  serverTxClientRx   +---------------------------------+
//                                        +---------------------+

// referenceChannel ...
type referenceChannel struct {
	refCount         int
	serverTxClientRx chan interface{}
	serverRxClientTx chan interface{}
}

// newReferenceChannel ...
func newReferenceChannel(size int) *referenceChannel {
	return &referenceChannel{
		refCount:         0,
		serverTxClientRx: make(chan interface{}, size),
		serverRxClientTx: make(chan interface{}, size),
	}
}

// deleteReferenceChannel ...
func deleteReferenceChannel(smc *referenceChannel) {
	if smc != nil {
		close(smc.serverTxClientRx)
		close(smc.serverRxClientTx)
	}
}

// ReferenceTransportConnection ...
type ReferenceTransportConnection struct {
	name      string
	forServer bool
	closed    bool
	channel   *referenceChannel
}

// NewReferenceTransportConnection ...
func NewReferenceTransportConnection(
	name string,
	forServer bool,
	channel *referenceChannel) TransportConnection {

	return &ReferenceTransportConnection{
		name:      name,
		forServer: forServer,
		closed:    false,
		channel:   channel,
	}
}

// GetName ...
func (c *ReferenceTransportConnection) GetName() string {
	return c.name
}

// Read ...
func (c *ReferenceTransportConnection) Read(p []byte) (n int, err error) {
	return 0, errors.New("ReferenceTransportConnection:Read: does not support this call")
}

// Write ...
func (c *ReferenceTransportConnection) Write(p []byte) (n int, err error) {
	return 0, errors.New("ReferenceTransportConnection:Write: does not support this call")
}

// UseReference ...
func (c *ReferenceTransportConnection) UseReference() bool {
	return true
}

// GetReference ...
func (c *ReferenceTransportConnection) GetReference() (p interface{}, err error) {
	if c.closed {
		fmt.Println("GetReference failed as connection is closed")
		return nil, errors.New("SetReference: connection is closed")
	}

	var data interface{}

	if c.forServer {
		data = <-c.channel.serverRxClientTx
	} else {
		data = <-c.channel.serverTxClientRx
	}

	return data, nil
}

// SetReference ...
func (c *ReferenceTransportConnection) SetReference(p interface{}) error {
	if c.closed {
		fmt.Println("SetReference failed as connection is closed")
		return errors.New("SetReference: connection is closed")
	}

	if p == nil {
		return errors.New("SetReference: null input")
	}

	if c.forServer {
		c.channel.serverTxClientRx <- p
	} else {
		c.channel.serverRxClientTx <- p
	}

	return nil
}

// Close ...
func (c *ReferenceTransportConnection) Close() {
	c.closed = true
}

// Closed ...
func (c *ReferenceTransportConnection) Closed() bool {
	return c.closed
}

// ReferenceTransport ...
type ReferenceTransport struct {
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

// NewReferenceTransport ...
func NewReferenceTransport(name string) Transport {
	return &ReferenceTransport{
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
func (sm *ReferenceTransport) parseOptions() {
	size, exists := OptionParseInt(sm.GetOption("MaxQueueSize"), 512)
	sm.queueSize = size
	if exists {
		fmt.Println("ReferenceTransport: option MaxQueueSize:", sm.queueSize)
	}
}

// ForServer ...
func (sm *ReferenceTransport) ForServer(forServer bool) Transport {
	sm.forServer = forServer
	return sm
}

// GetName ...
func (sm *ReferenceTransport) GetName() string {
	return sm.name
}

// SetOption ...
func (sm *ReferenceTransport) SetOption(name string, value interface{}) Transport {
	sm.options[name] = value
	return sm
}

// GetOption ...
func (sm *ReferenceTransport) GetOption(name string) interface{} {
	if value, ok := sm.options[name]; ok {
		return value
	}
	return nil
}

// SetAddress ...
func (sm *ReferenceTransport) SetAddress(address string) Transport {
	sm.address = address
	return sm
}

// GetAddress ...
func (sm *ReferenceTransport) GetAddress() string {
	return sm.address
}

// NextConnection ...
func (sm *ReferenceTransport) NextConnection() TransportConnection {
	next := <-sm.next
	sm.connection = next
	return sm.connection
}

// Two Transports(Client side and Server side) should use the same
// referenceChannel instance.
// Keep a global list for reusing
var mutex sync.Mutex
var chans map[string]*referenceChannel = make(map[string]*referenceChannel, 32)

// Run ...
func (sm *ReferenceTransport) Run() Transport {
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
		ch = newReferenceChannel(sm.queueSize)
		chans[sm.address] = ch
	}

	ch.refCount++

	sm.next <- NewReferenceTransportConnection(sm.name, sm.forServer, ch)

	mutex.Unlock()

	return sm
}

// Stop ...
func (sm *ReferenceTransport) Stop() Transport {
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
			deleteReferenceChannel(ch)
		}
	}

	mutex.Unlock()

	return sm
}
