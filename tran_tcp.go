// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// TCPTransportConnection ...
type TCPTransportConnection struct {
	name   string
	conn   net.Conn
	closed bool
}

// NewTCPTransportConnection ...
func NewTCPTransportConnection(name string, conn net.Conn) TransportConnection {
	return &TCPTransportConnection{
		name:   name,
		conn:   conn,
		closed: false,
	}
}

// GetName ...
func (c *TCPTransportConnection) GetName() string {
	return c.name
}

// Read ...
func (c *TCPTransportConnection) Read(p []byte) (n int, err error) {
	if c.closed {
		fmt.Println("read failed as connection", c.name, " is closed")
		return 0, nil
	}

	n, err = c.conn.Read(p)
	if err != nil {
		if err != io.EOF {
			fmt.Println("connection", c.name, "read failed:", err)
		}
	}

	return n, err
}

// Write ...
func (c *TCPTransportConnection) Write(p []byte) (n int, err error) {
	if c.closed {
		fmt.Println("write failed as connection", c.name, " is closed")
		return 0, nil
	}

	return c.conn.Write(p)
}

// UseReference ...
func (c *TCPTransportConnection) UseReference() bool {
	return false
}

// GetReference ...
func (c *TCPTransportConnection) GetReference() (p interface{}, err error) {
	return nil, errors.New("GetReference: Does not support this call")
}

// SetReference ...
func (c *TCPTransportConnection) SetReference(p interface{}) error {
	return errors.New("GetReference: Does not support this call")
}

// Close ...
func (c *TCPTransportConnection) Close() {
	c.closed = true
	c.conn.Close()
}

// Closed ...
func (c *TCPTransportConnection) Closed() bool {
	return c.closed
}

// TCPTransport ...
type TCPTransport struct {
	name    string
	options map[string]interface{}
	sync.Mutex
	address     string
	isRunning   bool
	forServer   bool
	connections []TransportConnection
	next        chan TransportConnection
	// for server
	listener net.Listener
	// for client
	maxRetryCount         int
	retryIntervalInSecond int
}

// NewTCPTransport ...
func NewTCPTransport(name string) Transport {
	return &TCPTransport{
		name:      name,
		options:   make(map[string]interface{}),
		isRunning: false,
		forServer: true,
		listener:  nil,
	}
}

// parseOptions ...
func (t *TCPTransport) parseOptions() {
	retry, exists := OptionParseInt(t.GetOption("MaxRetryCount"), 180)
	t.maxRetryCount = retry
	if exists {
		fmt.Println("TCPTransport: option MaxRetryCount:", t.maxRetryCount)
	}

	interval, exists := OptionParseInt(t.GetOption("RetryIntervalInSecond"), 1)
	t.retryIntervalInSecond = interval
	if exists {
		fmt.Println("TCPTransport: option RetryIntervalInSecond:", t.retryIntervalInSecond)
	}
}

// doInit ...
func (t *TCPTransport) doInit() {
	t.connections = make([]TransportConnection, 0)
	t.next = make(chan TransportConnection, 16)
}

// saveConnections ...
func (t *TCPTransport) saveConnection(tc TransportConnection) {
	t.Lock()
	defer t.Unlock()

	for _, c := range t.connections {
		if c == tc {
			return
		}
	}
	t.connections = append(t.connections, tc)
}

// dropConnections ...
func (t *TCPTransport) dropConnections() {
	for _, c := range t.connections {
		c.Close()
	}
}

// accept ...
func (t *TCPTransport) accept() {
	listener, err := net.Listen("tcp", t.address)
	if err != nil {
		log.Fatal(err)
	}

	t.listener = listener

	fmt.Println("Wait client connection ...")

	for {
		next, err := t.listener.Accept()
		if err != nil {
			break
		}

		t.next <- NewTCPTransportConnection(
			next.RemoteAddr().String(),
			next)
	}

	t.Stop()
}

// connect ...
func (t *TCPTransport) connect() {
	fmt.Println("Dial server ...")

	for i := 0; i < t.maxRetryCount; i++ {
		connection, err := net.DialTimeout("tcp", t.address, time.Second)

		if connection != nil && err == nil {
			t.next <- NewTCPTransportConnection(
				connection.RemoteAddr().String(),
				connection)
			return
		}

		fmt.Println(err, "retry", i)

		time.Sleep(time.Second * time.Duration(t.retryIntervalInSecond))
	}

	fmt.Println("Timeout to connect server")

	t.Stop()
}

// ForServer ...
func (t *TCPTransport) ForServer(forServer bool) Transport {
	t.forServer = forServer
	return t
}

// GetName ...
func (t *TCPTransport) GetName() string {
	return t.name
}

// SetOption ...
func (t *TCPTransport) SetOption(name string, value interface{}) Transport {
	t.options[name] = value
	return t
}

// GetOption ...
func (t *TCPTransport) GetOption(name string) interface{} {
	if value, ok := t.options[name]; ok {
		return value
	}
	return nil
}

// SetAddress ...
func (t *TCPTransport) SetAddress(address string) Transport {
	t.address = address
	return t
}

// GetAddress ...
func (t *TCPTransport) GetAddress() string {
	return t.address
}

// NextConnection ...
func (t *TCPTransport) NextConnection() TransportConnection {
	next := <-t.next
	t.saveConnection(next)
	return next
}

// Run ...
func (t *TCPTransport) Run() Transport {
	t.Lock()
	defer t.Unlock()

	if t.isRunning {
		return t
	}

	t.isRunning = true

	t.parseOptions()
	t.doInit()

	if t.forServer {
		go t.accept()
	} else {
		go t.connect()
	}
	return t
}

// Stop ...
func (t *TCPTransport) Stop() Transport {
	t.Lock()
	defer t.Unlock()

	if !t.isRunning {
		return t
	}

	if t.listener != nil {
		t.listener.Close()
		t.listener = nil
	}

	close(t.next)

	t.dropConnections()

	t.isRunning = false

	return t
}
