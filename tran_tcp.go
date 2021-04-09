// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import (
	"fmt"
	"io"
	"log"
	"net"
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
func (c *TCPTransportConnection) Write(p []byte) (n int, err error) {
	if c.closed {
		fmt.Println("connection is closed")
		return 0, nil
	}

	return c.conn.Write(p)
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
	name       string
	address    string
	forServer  bool
	connection chan TransportConnection
}

// NewTCPTransport ...
func NewTCPTransport(name string) Transport {
	return &TCPTransport{
		name:       name,
		forServer:  true,
		connection: make(chan TransportConnection, 16),
	}
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
	return <-t.connection
}

// Run ...
func (t *TCPTransport) Run() Transport {
	if t.forServer {
		go t.accept()
	} else {
		go t.connect()
	}
	return t
}

// accept ...
func (t *TCPTransport) accept() {
	listener, err := net.Listen("tcp", t.address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Wait client connection ...")

	for {
		connection, err := listener.Accept()
		if err == nil {
			t.connection <- NewTCPTransportConnection(
				connection.RemoteAddr().String(),
				connection)
		}
	}
}

// connect ...
func (t *TCPTransport) connect() {
	fmt.Println("Dial TCP Server ...")

	for {
		connection, err := net.DialTimeout("tcp", t.address, time.Second*10)
		if err == nil {
			t.connection <- NewTCPTransportConnection(
				connection.RemoteAddr().String(),
				connection)
			return
		}

		fmt.Println(err, ",retry ...")

		time.Sleep(time.Second)
	}
}
