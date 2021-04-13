// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// LowerDeck manages transports
type LowerDeck struct {
	Base
}

// NewLowerDeck returns a new instance
func NewLowerDeck() DataProcessor {
	ld := &LowerDeck{
		NewBaseInstance("LowerDeck"),
	}
	return ld.Base.SetWhere(ld)
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

// Run monitor new coming connection with routine
// and receive data from any new connection with routine
func (ld *LowerDeck) Run() DataProcessor {
	for _, transport := range ld.ustack.GetTransport() {
		tp := transport
		// New routinue to wait connections
		go func() {
			for {
				// this call will be blocked until new connection coming
				connection := tp.NextConnection()

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

	return ld
}
