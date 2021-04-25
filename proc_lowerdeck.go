// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// LowerDeck manages transports
type LowerDeck struct {
	ProcBase
}

// NewLowerDeck returns a new instance
func NewLowerDeck() DataProcessor {
	ld := &LowerDeck{
		NewProcBaseInstance("LowerDeck"),
	}
	return ld.ProcBase.SetWhere(ld)
}

func (ld *LowerDeck) closeConnection(c TransportConnection) {
	// close first
	c.Close()

	// publish event
	ld.ustack.PublishEvent(Event{
		Type:   UStackEventConnectionClosed,
		Source: ld,
		Data:   c,
	})
}

// acceptTransport ...
func (ld *LowerDeck) acceptTransport(tp Transport) {

	tp.Run()

	// New routinue to wait connections
	go func() {
		for {
			// this call will be blocked until new connection coming
			connection := tp.NextConnection()

			if connection == nil {
				return
			}

			fmt.Println("New connection:", connection.GetName(), "on transport:", tp.GetName())

			// publish event
			ld.ustack.PublishEvent(Event{
				Type:   UStackEventNewConnection,
				Source: ld,
				Data:   connection,
			})

			// New routine to continue receive data from connection
			go func() {
				for {
					ub := UBufAlloc(ld.ustack.GetMTU())

					n, err := ub.ReadFrom(connection)
					if n == 0 || err != nil {
						ld.closeConnection(connection)
						return
					}

					// invoke the uplayer
					ld.upper.OnLowerData(
						NewUStackContext().
							SetConnection(connection).
							SetBuffer(ub))
				}
			}()
		}
	}()
}

// deleteTransport ...
func (ld *LowerDeck) deleteTransport(tp Transport) {
	tp.Stop()
}

// OnUpperData sends ulayer data with connection
func (ld *LowerDeck) OnUpperData(context Context) {
	ub := context.GetBuffer()
	if ub == nil {
		return
	}

	connection := context.GetConnection()
	if connection == nil {
		return
	}

	_, err := ub.WriteTo(connection)
	if err != nil {
		fmt.Printf("Connection is closed: %s\n", connection.GetName())
		ld.closeConnection(connection)
	}
}

// OnEvent is called when any event hanppen
func (ld *LowerDeck) OnEvent(event Event) {
	tp, ok := event.Data.(Transport)
	if !ok {
		return
	}

	if event.Type == UStackEventTransportAdded {
		ld.acceptTransport(tp)
	} else if event.Type == UStackEventTransportDeleted {
		ld.deleteTransport(tp)
	}
}

// Run monitor new coming connection with routine
// and receive data from any new connection with routine
func (ld *LowerDeck) Run() DataProcessor {
	for _, tp := range ld.ustack.GetTransport() {
		ld.acceptTransport(tp)
	}

	return ld
}
