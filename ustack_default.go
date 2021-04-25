// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

const (
	defaultMTU int = 2048
)

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
func (c *DefaultUStackContext) GetOption(name string) interface{} {
	if value, ok := c.options[name]; ok {
		return value
	}
	return nil
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
	mtu        int
	options    map[string]interface{}
	features   []Feature
	endpoints  []EndPoint
	transports []Transport
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
		mtu:        defaultMTU,
		options:    make(map[string]interface{}),
		features:   nil,
		endpoints:  nil,
		transports: nil,
		upperDeck:  nil,
		processors: nil,
		overhead:   0,
		lowerDeck:  nil,
		listeners:  nil,
	}
}

// parseOptions ...
func (u *DefaultUStack) parseOptions() {
	mtu, exists := OptionParseInt(u.GetOption("MTU"), defaultMTU)
	u.mtu = mtu
	if exists {
		fmt.Println("UStack: option MTU:", u.mtu)
	}
}

// build ...
func (u *DefaultUStack) build() {
	u.parseOptions()

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
		u.processors[i].SetUStack(u)
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

// AddFeature ...
func (u *DefaultUStack) AddFeature(feature Feature) UStack {
	u.features = append(u.features, feature)
	return u
}

// GetFeature ...
func (u *DefaultUStack) GetFeature(name string) Feature {
	for _, feature := range u.features {
		if feature.GetName() == name {
			return feature
		}
	}
	return nil
}

// GetFeatures ...
func (u *DefaultUStack) GetFeatures() []Feature {
	return u.features
}

// AddEndPoint ...
func (u *DefaultUStack) AddEndPoint(ep EndPoint) UStack {
	for _, endpoint := range u.endpoints {
		if endpoint == ep {
			// was added
			return u
		}
	}

	u.PublishEvent(Event{
		Type:   UStackEventEndpointAdded,
		Source: u,
		Data:   ep,
	})

	u.endpoints = append(u.endpoints, ep)
	return u
}

// DeleteEndPoint ...
func (u *DefaultUStack) DeleteEndPoint(ep EndPoint) UStack {
	for i, endpoint := range u.endpoints {
		if endpoint != ep {
			continue
		}

		u.PublishEvent(Event{
			Type:   UStackEventEndpointDeleted,
			Source: u,
			Data:   ep,
		})

		// delete
		u.endpoints = append(u.endpoints[:i], u.endpoints[i+1:]...)

		break
	}

	return u
}

// GetEndPoint ...
func (u *DefaultUStack) GetEndPoint() []EndPoint {
	return u.endpoints
}

// AppendDataProcessor ...
func (u *DefaultUStack) AppendDataProcessor(dp DataProcessor) UStack {
	u.processors = append(u.processors, dp)
	return u
}

// GetOverhead returns all data processors overhead
func (u *DefaultUStack) GetOverhead() int {
	return u.overhead
}

// GetMTU returns maximum transmission unit size
func (u *DefaultUStack) GetMTU() int {
	return u.mtu
}

// AddTransport ...
func (u *DefaultUStack) AddTransport(tp Transport) UStack {
	for _, transport := range u.transports {
		if transport == tp {
			// was added
			return u
		}
	}

	u.PublishEvent(Event{
		Type:   UStackEventTransportAdded,
		Source: u,
		Data:   tp,
	})

	u.transports = append(u.transports, tp)
	return u
}

// DeleteTransport ...
func (u *DefaultUStack) DeleteTransport(tp Transport) UStack {
	for i, transport := range u.transports {
		if transport != tp {
			continue
		}

		// delete
		u.transports = append(u.transports[:i], u.transports[i+1:]...)

		u.PublishEvent(Event{
			Type:   UStackEventTransportDeleted,
			Source: u,
			Data:   tp,
		})

		break
	}

	return u
}

// GetTransport ...
func (u *DefaultUStack) GetTransport() []Transport {
	return u.transports
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
	if u.upperDeck != nil {
		u.upperDeck.OnEvent(event)
	}
	for _, processor := range u.processors {
		processor.OnEvent(event)
	}
	if u.lowerDeck != nil {
		u.lowerDeck.OnEvent(event)
	}
	return u
}

// Run ...
func (u *DefaultUStack) Run() UStack {
	u.build()

	for _, ft := range u.features {
		ft.Run()
	}

	u.upperDeck.Run()

	for _, dp := range u.processors {
		dp.Run()
	}

	u.lowerDeck.Run()

	return u
}
