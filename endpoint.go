// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// EndPointData ...
type EndPointData interface {
	GetConnection() TransportConnection
	GetData() interface{}
}

// EndPoint ...
type EndPoint interface {
	GetName() string
	GetSession() int
	SetDataListener(listener func(EndPoint, EndPointData)) EndPoint
	GetTxChannel() chan EndPointData
	GetRxChannel() chan EndPointData
	SetEventListener(listener func(EndPoint, Event)) EndPoint
	OnEvent(event Event)
}
