// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// DefaultEndPointData ...
type DefaultEndPointData struct {
	connection            TransportConnection
	data                  interface{}
	hasDestinationSession bool
	destinationSession    int
}

// NewEndPointData ...
func NewEndPointData() EndPointData {
	return &DefaultEndPointData{
		connection:            nil,
		data:                  nil,
		hasDestinationSession: false,
		destinationSession:    0,
	}
}

// GetConnection ...
func (epd *DefaultEndPointData) SetConnection(c TransportConnection) EndPointData {
	epd.connection = c
	return epd
}

// GetConnection ...
func (epd *DefaultEndPointData) GetConnection() TransportConnection {
	return epd.connection
}

// SetData ...
func (epd *DefaultEndPointData) SetData(data interface{}) EndPointData {
	epd.data = data
	return epd
}

// GetData ...
func (epd *DefaultEndPointData) GetData() interface{} {
	return epd.data
}

// HasDestinationSession ...
func (epd *DefaultEndPointData) HasDestinationSession() bool {
	return epd.hasDestinationSession
}

// HasDestinationSession ...
func (epd *DefaultEndPointData) SetDestinationSession(session int) EndPointData {
	epd.destinationSession = session
	epd.hasDestinationSession = true
	return epd
}

// GetDestinationSession ...
func (epd *DefaultEndPointData) GetDestinationSession() int {
	return epd.destinationSession
}

// HasDestinationSession ...
func (epd *DefaultEndPointData) ClearDestinationSession() EndPointData {
	epd.hasDestinationSession = false
	return epd
}

// DefaultEndPoint ...
type DefaultEndPoint struct {
	name          string
	session       int
	txChannel     chan EndPointData
	rxChannel     chan EndPointData
	eventListener func(EndPoint, Event)
	dataListener  func(EndPoint, EndPointData)
	inAutoReceive bool
}

// autoReceive ...
func (ep *DefaultEndPoint) autoReceive() {
	if ep.dataListener != nil {
		ep.inAutoReceive = true
		for epd := range ep.rxChannel {
			if ep.dataListener == nil {
				ep.inAutoReceive = false
				return
			}
			ep.dataListener(ep, epd)
		}
	}
}

// NewEndPoint ...
func NewEndPoint(name string, session int) EndPoint {
	return &DefaultEndPoint{
		name:          name,
		session:       session,
		txChannel:     make(chan EndPointData, 512),
		rxChannel:     make(chan EndPointData, 512),
		eventListener: nil,
		dataListener:  nil,
		inAutoReceive: false,
	}
}

// SetEventListener ...
func (ep *DefaultEndPoint) SetEventListener(listener func(EndPoint, Event)) EndPoint {
	ep.eventListener = listener
	return ep
}

// SetDataListener ...
func (ep *DefaultEndPoint) SetDataListener(listener func(EndPoint, EndPointData)) EndPoint {
	ep.dataListener = listener
	if !ep.inAutoReceive {
		go ep.autoReceive()
	}
	return ep
}

// SetName ...
func (ep *DefaultEndPoint) SetName(name string) EndPoint {
	ep.name = name
	return ep
}

// GetName ...
func (ep *DefaultEndPoint) GetName() string {
	return ep.name
}

// SetName ...
func (ep *DefaultEndPoint) SetSession(session int) EndPoint {
	ep.session = session
	return ep
}

// GetPort ...
func (ep *DefaultEndPoint) GetSession() int {
	return ep.session
}

// GetTxChannel ...
func (ep *DefaultEndPoint) GetTxChannel() chan EndPointData {
	return ep.txChannel
}

// GetRxChannel ...
func (ep *DefaultEndPoint) GetRxChannel() chan EndPointData {
	return ep.rxChannel
}

// OnEvent ...
func (ep *DefaultEndPoint) OnEvent(event Event) {
	if ep.eventListener != nil {
		ep.eventListener(ep, event)
	}
}
