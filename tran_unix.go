// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
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
		fmt.Println("connection is closed")
		return 0, nil
	}

	n, err = c.conn.Read(p)
	if err != nil {
		if err != io.EOF {
			fmt.Println("read error:", err)
		}
	}

	return n, err
}

// Write ...
func (c *UDSTransportConnection) Write(p []byte) (n int, err error) {
	if c.closed {
		fmt.Println("connection is closed")
		return 0, nil
	}

	return c.conn.Write(p)
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
	name       string
	filename   string
	forServer  bool
	connection chan TransportConnection
}

// NewUDSTransport ...
func NewUDSTransport(name string) Transport {
	return &UDSTransport{
		name:       name,
		forServer:  true,
		connection: make(chan TransportConnection, 16),
	}
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
	return <-uds.connection
}

// Run ...
func (uds *UDSTransport) Run() Transport {
	if uds.forServer {
		go uds.accept()
	} else {
		go uds.connect()
	}
	return uds
}

// accept ...
func (uds *UDSTransport) accept() {
	os.Remove(uds.filename)

	listener, err := net.Listen("unix", uds.filename)
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(uds.filename)
	defer listener.Close()

	fmt.Println("Wait client connection ...")

	for {
		connection, err := listener.Accept()
		if err == nil {
			uds.connection <- NewUDSTransportConnection(
				connection.RemoteAddr().String(),
				connection)
		}
	}
}

// connect ...
func (uds *UDSTransport) connect() {
	fmt.Println("Dial UDS Server ...")

	for {
		connection, err := net.DialTimeout("unix", uds.filename, time.Second*10)
		if err == nil {
			uds.connection <- NewUDSTransportConnection(
				connection.RemoteAddr().String(),
				connection)
			return
		}

		fmt.Println(err, ",retry ...")

		time.Sleep(time.Second)
	}
}
