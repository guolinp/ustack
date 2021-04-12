// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// UpperDeck ...
type UpperDeck struct {
	Base
	endpoints []EndPoint
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
func NewUpperDeck() DataProcessor {
	ud := &UpperDeck{
		Base:      NewBaseInstance("UpperDeck"),
		endpoints: nil,
	}
	return ud.Base.SetWhere(ud)
}

// OnLowerData ...
func (ud *UpperDeck) OnLowerData(context Context) {
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

// Run ...
func (ud *UpperDeck) Run() DataProcessor {
	ud.build()

	ldp := ud.lower
	for _, ep := range ud.endpoints {
		session := ep.GetSession()
		txchan := ep.GetTxChannel()
		go func() {
			for epd := range txchan {
				ldp.OnUpperData(
					NewUStackContext().
						SetConnection(epd.GetConnection()).
						SetOption("session", session).
						SetOption("message", epd.GetData()))
			}
		}()
	}

	return ud
}
