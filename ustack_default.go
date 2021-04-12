// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// DefaultUStackContext ...
type DefaultUStackContext struct {
	connection TransportConnection
	ubuf       *UBuf
	options    map[string]interface{}
}

// NewUStackContext ...
func NewUStackContext() Context {
	return &DefaultUStackContext{
		connection: nil,
		ubuf:       nil,
		options:    make(map[string]interface{}),
	}
}

// SetConnection ...
func (c *DefaultUStackContext) SetConnection(connection TransportConnection) Context {
	c.connection = connection
	return c
}

// GetConnection ...
func (c *DefaultUStackContext) GetConnection() TransportConnection {
	return c.connection
}

// SetOption ...
func (c *DefaultUStackContext) SetOption(name string, value interface{}) Context {
	c.options[name] = value
	return c
}

// GetOption ...
func (c *DefaultUStackContext) GetOption(name string) (interface{}, bool) {
	if value, ok := c.options[name]; ok {
		return value, true
	}
	return nil, false
}

// SetCacheData ...
func (c *DefaultUStackContext) SetBuffer(ubuf *UBuf) Context {
	c.ubuf = ubuf
	return c
}

// GetCacheData ...
func (c *DefaultUStackContext) GetBuffer() *UBuf {
	return c.ubuf
}

// DefaultUStack ...
type DefaultUStack struct {
	name       string
	options    map[string]interface{}
	endpoints  []EndPoint
	transport  Transport
	upperDeck  DataProcessor
	processors []DataProcessor
	overhead   int
	lowerDeck  DataProcessor
	listeners  []func(Event)
}

// NewUStack ...
func NewUStack() UStack {
	return &DefaultUStack{
		name:       "UStack",
		options:    make(map[string]interface{}),
		endpoints:  nil,
		transport:  nil,
		upperDeck:  nil,
		processors: nil,
		overhead:   0,
		lowerDeck:  nil,
		listeners:  nil,
	}
}

// build ...
func (u *DefaultUStack) build() {
	u.upperDeck = NewUpperDeck().SetUStack(u)
	u.lowerDeck = NewLowerDeck().SetUStack(u)

	count := len(u.processors)

	if count == 0 {
		u.upperDeck.SetLower(u.lowerDeck)
		u.lowerDeck.SetUpper(u.upperDeck)
	} else {
		u.upperDeck.SetLower(u.processors[0])
		u.processors[0].SetUpper(u.upperDeck)
		for i := 0; i < count-1; i++ {
			u.processors[i].SetLower(u.processors[i+1])
			u.processors[i+1].SetUpper(u.processors[i])
		}
		u.processors[count-1].SetLower(u.lowerDeck)
		u.lowerDeck.SetUpper(u.processors[count-1])
	}

	for i := 0; i < count; i++ {
		u.overhead += u.processors[i].GetOverhead()
	}
}

// SetName ...
func (u *DefaultUStack) SetName(name string) UStack {
	u.name = name
	return u
}

// GetName ...
func (u *DefaultUStack) GetName() string {
	return u.name
}

// SetOption ...
func (u *DefaultUStack) SetOption(name string, value interface{}) UStack {
	u.options[name] = value
	return u
}

// GetOption ...
func (u *DefaultUStack) GetOption(name string) interface{} {
	if value, ok := u.options[name]; ok {
		return value
	}
	return nil
}

// SetEndPoint ...
func (u *DefaultUStack) SetEndPoint(ep EndPoint) UStack {
	u.endpoints = append(u.endpoints, ep)
	return u
}

// GetEndPoint ...
func (u *DefaultUStack) GetEndPoint() []EndPoint {
	return u.endpoints
}

// SetDataProcessor ...
func (u *DefaultUStack) SetDataProcessor(dp DataProcessor) UStack {
	u.processors = append(u.processors, dp)
	return u
}

// GetOverhead returns all data processors overhead
func (u *DefaultUStack) GetOverhead() int {
	return u.overhead
}

// GetMTU returns maximum transmission unit size
func (u *DefaultUStack) GetMTU() int {
	// will get from options or Transport
	return 2048
}

// SetTransport ...
func (u *DefaultUStack) SetTransport(tp Transport) UStack {
	u.transport = tp
	return u
}

// GetTransport ...
func (u *DefaultUStack) GetTransport() Transport {
	return u.transport
}

// SetEventListener ...
func (u *DefaultUStack) SetEventListener(listener func(Event)) UStack {
	u.listeners = append(u.listeners, listener)
	return u
}

// PublishEvent ...
func (u *DefaultUStack) PublishEvent(event Event) UStack {
	for _, listener := range u.listeners {
		listener(event)
	}
	for _, endpoint := range u.endpoints {
		endpoint.OnEvent(event)
	}
	u.upperDeck.OnEvent(event)
	for _, processor := range u.processors {
		processor.OnEvent(event)
	}
	u.lowerDeck.OnEvent(event)
	return u
}

// Run ...
func (u *DefaultUStack) Run() UStack {
	u.build()
	u.upperDeck.Run()
	for _, dp := range u.processors {
		dp.Run()
	}
	u.lowerDeck.Run()
	u.transport.Run()
	return u
}
