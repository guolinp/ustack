// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// UpperDeck manages endpoints
type UpperDeck struct {
	ProcBase
	endpoints []EndPoint
}

func (ud *UpperDeck) build() {
	ud.endpoints = ud.ustack.GetEndPoint()
}

func (ud *UpperDeck) findEndPoint(session byte) EndPoint {
	for _, ep := range ud.endpoints {
		if ep.GetSession() == session {
			return ep
		}
	}
	return nil
}

// NewUpperDeck returns a new instance
func NewUpperDeck() DataProcessor {
	ud := &UpperDeck{
		ProcBase:      NewProcBaseInstance("UpperDeck"),
		endpoints: nil,
	}
	return ud.ProcBase.SetWhere(ud)
}

// OnLowerData finds the endpoint with session and pass data
func (ud *UpperDeck) OnLowerData(context Context) {
	message := context.GetOption("message")
	if message == nil {
		return
	}

	session, _ := OptionParseByte(context.GetOption("session"), 0)

	ep := ud.findEndPoint(session)
	if ep != nil {
		ep.GetRxChannel() <- NewEndPointData(context.GetConnection(), message)
	}
}

// Run ...
func (ud *UpperDeck) Run() DataProcessor {
	// do build
	ud.build()

	lower := ud.lower
	for _, ep := range ud.endpoints {
		session := ep.GetSession()
		txchan := ep.GetTxChannel()

		// create routinue for each endpoint
		// read data from endpoint and pass it to lowlayer
		go func() {
			for epd := range txchan {
				lower.OnUpperData(
					NewUStackContext().
						SetConnection(epd.GetConnection()).
						SetOption("session", session).
						SetOption("message", epd.GetData()))
			}
		}()
	}

	return ud
}
