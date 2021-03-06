// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// Context ...
type Context interface {
	SetConnection(connection TransportConnection) Context
	GetConnection() TransportConnection
	SetOption(name string, value interface{}) Context
	GetOption(name string) interface{}
	SetBuffer(ub *UBuf) Context
	GetBuffer() *UBuf
	SetMessage(message interface{}) Context
	GetMessage() interface{}
}

// UStack ...
type UStack interface {
	SetName(name string) UStack
	GetName() string

	SetOption(name string, value interface{}) UStack
	GetOption(name string) interface{}

	AddFeature(feature Feature) UStack
	GetFeature(name string) Feature
	GetFeatures() []Feature

	AddEndPoint(ep EndPoint) UStack
	DeleteEndPoint(ep EndPoint) UStack
	GetEndPoint() []EndPoint

	AppendDataProcessor(dp DataProcessor) UStack

	GetOverhead() int
	GetMTU() int

	AddTransport(tp Transport) UStack
	DeleteTransport(tp Transport) UStack
	GetTransport() []Transport

	SetEventListener(listener func(Event)) UStack
	PublishEvent(event Event) UStack

	Run() UStack
}
