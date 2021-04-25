// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import (
	"sync"
)

// UpperDeck manages endpoints
type UpperDeck struct {
	ProcBase
	sync.Mutex
	endpoints map[EndPoint]chan bool
}

func (ud *UpperDeck) existsEndpoint(ep EndPoint) bool {
	ud.Lock()
	defer ud.Unlock()

	for endpoint := range ud.endpoints {
		if endpoint == ep || endpoint.GetSession() == ep.GetSession() {
			return true
		}
	}
	return false
}

// accpetEndpoint ...
func (ud *UpperDeck) acceptEndpoint(ep EndPoint) {
	if ud.existsEndpoint(ep) {
		return
	}

	ud.Lock()
	defer ud.Unlock()

	ud.endpoints[ep] = make(chan bool, 1)

	session := ep.GetSession()
	txchan := ep.GetTxChannel()
	lower := ud.lower

	// create routinue for each endpoint
	// read data from endpoint and pass it to lower
	go func() {
		for {
			select {
			case stop := <-ud.endpoints[ep]:
				if stop {
					return
				}
			case epd := <-txchan:
				lower.OnUpperData(
					NewUStackContext().
						SetConnection(epd.GetConnection()).
						SetOption("session", session).
						SetOption("message", epd.GetData()))
			}
		}
	}()
}

// deleteTransport ...
func (ud *UpperDeck) deleteEndpoint(ep EndPoint) {
	ud.Lock()
	defer ud.Unlock()

	ch, ok := ud.endpoints[ep]
	if ok {
		ch <- true
		close(ch)
		delete(ud.endpoints, ep)
	}
}

func (ud *UpperDeck) findEndPoint(session int) EndPoint {
	ud.Lock()
	defer ud.Unlock()

	for ep := range ud.endpoints {
		if ep.GetSession() == session {
			return ep
		}
	}
	return nil
}

// NewUpperDeck returns a new instance
func NewUpperDeck() DataProcessor {
	ud := &UpperDeck{
		ProcBase:  NewProcBaseInstance("UpperDeck"),
		endpoints: make(map[EndPoint]chan bool, 1),
	}
	return ud.ProcBase.SetWhere(ud)
}

// OnLowerData finds the endpoint with session and pass data
func (ud *UpperDeck) OnLowerData(context Context) {
	message := context.GetOption("message")
	if message == nil {
		return
	}

	// the default is 0 if seesion is not enabled
	session, _ := OptionParseInt(context.GetOption("session"), 0)

	ep := ud.findEndPoint(session)
	if ep != nil {
		ep.GetRxChannel() <- NewEndPointData(context.GetConnection(), message)
	}
}

// OnEvent is called when any event hanppen
func (ud *UpperDeck) OnEvent(event Event) {
	ep, ok := event.Data.(EndPoint)
	if !ok {
		return
	}

	if event.Type == UStackEventEndpointAdded {
		ud.acceptEndpoint(ep)
	} else if event.Type == UStackEventEndpointDeleted {
		ud.deleteEndpoint(ep)
	}
}

// Run ...
func (ud *UpperDeck) Run() DataProcessor {
	for _, ep := range ud.ustack.GetEndPoint() {
		ud.acceptEndpoint(ep)
	}

	return ud
}
