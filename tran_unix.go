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
	"os"
	"sync"
	"time"
)

// UDSTransportConnection ...
type UDSTransportConnection struct {
	name   string
	conn   net.Conn
	closed bool
}

// NewUDSTransportConnection ...
func NewUDSTransportConnection(name string, conn net.Conn) TransportConnection {
	return &UDSTransportConnection{
		name:   name,
		conn:   conn,
		closed: false,
	}
}

// GetName ...
func (c *UDSTransportConnection) GetName() string {
	return c.name
}

// Read ...
func (c *UDSTransportConnection) Read(p []byte) (n int, err error) {
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
func (c *UDSTransportConnection) Write(p []byte) (n int, err error) {
	if c.closed {
		fmt.Println("write failed as connection", c.name, " is closed")
		return 0, nil
	}

	return c.conn.Write(p)
}

// UseReference ...
func (c *UDSTransportConnection) UseReference() bool {
	return false
}

// GetReference ...
func (c *UDSTransportConnection) GetReference() (p interface{}, err error) {
	return nil, errors.New("UDSTransportConnection:GetReference: does not support this call")
}

// SetReference ...
func (c *UDSTransportConnection) SetReference(p interface{}) error {
	return errors.New("UDSTransportConnection:SetReference: does not support this call")
}

// Close ...
func (c *UDSTransportConnection) Close() {
	c.closed = true
	c.conn.Close()
}

// Closed ...
func (c *UDSTransportConnection) Closed() bool {
	return c.closed
}

// UDSTransport ...
type UDSTransport struct {
	name    string
	options map[string]interface{}
	sync.Mutex
	filename    string
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

// NewUDSTransport ...
func NewUDSTransport(name string) Transport {
	return &UDSTransport{
		name:      name,
		options:   make(map[string]interface{}),
		isRunning: false,
		forServer: true,
		listener:  nil,
	}
}

// parseOptions ...
func (uds *UDSTransport) parseOptions() {
	retry, exists := OptionParseInt(uds.GetOption("MaxRetryCount"), 180)
	uds.maxRetryCount = retry
	if exists {
		fmt.Println("UDSTransport: option MaxRetryCount:", uds.maxRetryCount)
	}

	interval, exists := OptionParseInt(uds.GetOption("RetryIntervalInSecond"), 1)
	uds.retryIntervalInSecond = interval
	if exists {
		fmt.Println("UDSTransport: option RetryIntervalInSecond:", uds.retryIntervalInSecond)
	}
}

// doInit ...
func (uds *UDSTransport) doInit() {
	uds.connections = make([]TransportConnection, 0)
	uds.next = make(chan TransportConnection, 16)
}

// saveConnections ...
func (uds *UDSTransport) saveConnection(tc TransportConnection) {
	if tc == nil {
		return
	}

	uds.Lock()
	defer uds.Unlock()

	for _, c := range uds.connections {
		if c == tc {
			return
		}
	}
	uds.connections = append(uds.connections, tc)
}

// dropConnections ...
func (uds *UDSTransport) dropConnections() {
	for _, c := range uds.connections {
		c.Close()
	}
}

// accept ...
func (uds *UDSTransport) accept() {
	os.Remove(uds.filename)

	listener, err := net.Listen("unix", uds.filename)
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(uds.filename)

	uds.listener = listener

	fmt.Println("Wait client connection ...")

	for {
		next, err := uds.listener.Accept()
		if err != nil {
			break
		}

		uds.next <- NewUDSTransportConnection(
			next.RemoteAddr().String(),
			next)
	}

	uds.Stop()
}

// connect ...
func (uds *UDSTransport) connect() {
	fmt.Println("Dial server ...")

	for i := 0; i < uds.maxRetryCount; i++ {
		connection, err := net.DialTimeout("unix", uds.filename, time.Second)

		if connection != nil && err == nil {
			uds.next <- NewUDSTransportConnection(
				connection.RemoteAddr().String(),
				connection)
			return
		}

		fmt.Println(err, "retry", i+1)

		time.Sleep(time.Second * time.Duration(uds.retryIntervalInSecond))
	}

	fmt.Println("Timeout to connect server")
	uds.Stop()
}

// ForServer ...
func (uds *UDSTransport) ForServer(forServer bool) Transport {
	uds.forServer = forServer
	return uds
}

// GetName ...
func (uds *UDSTransport) GetName() string {
	return uds.name
}

// SetOption ...
func (uds *UDSTransport) SetOption(name string, value interface{}) Transport {
	uds.options[name] = value
	return uds
}

// GetOption ...
func (uds *UDSTransport) GetOption(name string) interface{} {
	if value, ok := uds.options[name]; ok {
		return value
	}
	return nil
}

// SetAddress ...
func (uds *UDSTransport) SetAddress(address string) Transport {
	uds.filename = address
	return uds
}

// GetAddress ...
func (uds *UDSTransport) GetAddress() string {
	return uds.filename
}

// NextConnection ...
func (uds *UDSTransport) NextConnection() TransportConnection {
	next := <-uds.next
	uds.saveConnection(next)
	return next
}

// Run ...
func (uds *UDSTransport) Run() Transport {
	uds.Lock()
	defer uds.Unlock()

	if uds.isRunning {
		return uds
	}

	uds.isRunning = true

	uds.parseOptions()
	uds.doInit()

	if uds.forServer {
		go uds.accept()
	} else {
		go uds.connect()
	}
	return uds
}

// Stop ...
func (uds *UDSTransport) Stop() Transport {
	uds.Lock()
	defer uds.Unlock()

	if !uds.isRunning {
		return uds
	}

	if uds.listener != nil {
		uds.listener.Close()
		uds.listener = nil
	}

	close(uds.next)

	uds.dropConnections()

	uds.isRunning = false

	return uds
}
