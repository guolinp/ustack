// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// EndPointData ...
type EndPointData interface {
	SetConnection(c TransportConnection) EndPointData
	GetConnection() TransportConnection
	SetData(data interface{}) EndPointData
	GetData() interface{}
	HasDestinationSession() bool
	SetDestinationSession(session int) EndPointData
	GetDestinationSession() int
	ClearDestinationSession() EndPointData
}

// EndPoint ...
type EndPoint interface {
	SetName(name string) EndPoint
	GetName() string
	SetSession(session int) EndPoint
	GetSession() int
	GetTxChannel() chan EndPointData
	GetRxChannel() chan EndPointData
	SetDataListener(listener func(EndPoint, EndPointData)) EndPoint
	SetEventListener(listener func(EndPoint, Event)) EndPoint
	OnEvent(event Event)
}
