// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// TransportConnection ...
type TransportConnection interface {
	GetName() string
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close()
	Closed() bool
}

// Transport ...
type Transport interface {
	GetName() string
	SetOption(name string, value interface{}) Transport
	GetOption(name string) interface{}
	ForServer(bool) Transport
	SetAddress(address string) Transport
	GetAddress() string
	NextConnection() TransportConnection
	Run() Transport
	Stop() Transport
}
