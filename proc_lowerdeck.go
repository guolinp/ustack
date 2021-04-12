// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// LowerDeck ...
type LowerDeck struct {
	Base
}

// NewLowerDeck ...
func NewLowerDeck() DataProcessor {
	ld := &LowerDeck{
		NewBaseInstance("LowerDeck"),
	}
	return ld.Base.SetWhere(ld)
}

func (ld *LowerDeck) closeConnection(c TransportConnection) {
	c.Close()
	ld.ustack.PublishEvent(Event{
		Type:   UStackEventConnectionClosed,
		Source: ld,
		Data:   c,
	})
}

// OnUpperData ...
func (ld *LowerDeck) OnUpperData(context Context) {
	connection := context.GetConnection()
	if connection == nil {
		return
	}

	ub := context.GetBuffer()
	if ub == nil {
		return
	}

	_, err := ub.WriteTo(connection)
	if err != nil {
		fmt.Printf("Connection is closed: %s\n", connection.GetName())
		ld.closeConnection(connection)
	}
}

// Run ...
func (ld *LowerDeck) Run() DataProcessor {
	tp := ld.ustack.GetTransport()

	go func() {
		for {
			connection := tp.NextConnection()

			fmt.Println("New connection:", connection.GetName())

			ld.ustack.PublishEvent(Event{
				Type:   UStackEventNewConnection,
				Source: ld,
				Data:   connection,
			})

			go func() {
				for {
					ub := UBufAlloc(4096)

					n, err := ub.ReadFrom(connection)
					if n == 0 || err != nil {
						ld.closeConnection(connection)
						return
					}

					ld.upper.OnLowerData(
						NewUStackContext().
							SetConnection(connection).
							SetBuffer(ub))
				}
			}()
		}
	}()

	return ld
}
