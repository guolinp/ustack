// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// UpperDeck ...
type UpperDeck struct {
	name               string
	ustack             UStack
	options            map[string]interface{}
	upperDataProcessor DataProcessor
	lowerDataProcessor DataProcessor
	endpoints          []EndPoint
}

// build ...
func (ud *UpperDeck) build() {
	ud.endpoints = ud.ustack.GetEndPoint()
}

// findEndPoint ...
func (ud *UpperDeck) findEndPoint(session byte) EndPoint {
	for _, ep := range ud.endpoints {
		if ep.GetSession() == session {
			return ep
		}
	}
	return nil
}

// NewUpperDeck ...
func NewUpperDeck() *UpperDeck {
	return &UpperDeck{
		name:    "UpperDeck",
		options: make(map[string]interface{}),
	}
}

// SetName ...
func (ud *UpperDeck) SetName(name string) DataProcessor {
	ud.name = name
	return ud
}

// GetName ...
func (ud *UpperDeck) GetName() string {
	return ud.name
}

// SetOption ...
func (ud *UpperDeck) SetOption(name string, value interface{}) DataProcessor {
	ud.options[name] = value
	return ud
}

// GetOption ...
func (ud *UpperDeck) GetOption(name string) interface{} {
	if value, ok := ud.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable ...
func (ud *UpperDeck) SetEnable(enable bool) DataProcessor {
	return ud
}

// ForServer ...
func (ud *UpperDeck) ForServer(forServer bool) DataProcessor {
	return ud
}

// SetUStack ...
func (ud *UpperDeck) SetUStack(u UStack) DataProcessor {
	ud.ustack = u
	return ud
}

// SetUpperDataProcessor ...
func (ud *UpperDeck) SetUpperDataProcessor(dp DataProcessor) DataProcessor {
	ud.upperDataProcessor = dp
	return ud
}

// SetLowerDataProcessor ...
func (ud *UpperDeck) SetLowerDataProcessor(dp DataProcessor) DataProcessor {
	ud.lowerDataProcessor = dp
	return ud
}

// OnUpperPush ...
func (ud *UpperDeck) OnUpperPush(context Context) {
}

// OnLowerPush ...
func (ud *UpperDeck) OnLowerPush(context Context) {
	var session byte = 0
	sessionItf, ok := context.GetOption("session")
	if ok {
		sessionByte, ok := sessionItf.(byte)
		if ok {
			session = sessionByte
		}
	}

	ep := ud.findEndPoint(session)
	if ep == nil {
		return
	}

	message, ok := context.GetOption("message")
	if !ok {
		return
	}

	ep.GetRxChannel() <- NewEndPointData(context.GetConnection(), message)
}

// OnEvent ...
func (ud *UpperDeck) OnEvent(event Event) {
}

// Run ...
func (ud *UpperDeck) Run() DataProcessor {
	ud.build()

	ldp := ud.lowerDataProcessor
	for _, ep := range ud.endpoints {
		session := ep.GetSession()
		txchan := ep.GetTxChannel()
		go func() {
			for epd := range txchan {
				ldp.OnUpperPush(
					NewUStackContext().
						SetConnection(epd.GetConnection()).
						SetOption("session", session).
						SetOption("message", epd.GetData()))
			}
		}()
	}

	return ud
}
