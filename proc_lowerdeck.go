// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// LowerDeck ...
type LowerDeck struct {
	name               string
	ustack             UStack
	options            map[string]interface{}
	upperDataProcessor DataProcessor
	lowerDataProcessor DataProcessor
}

// SetName ...
func (ld *LowerDeck) closeConnection(c TransportConnection) {
	c.Close()
	ld.ustack.PublishEvent(Event{
		Type:   UStackEventConnectionClosed,
		Source: ld,
		Data:   c,
	})
}

// NewLowerDeck ...
func NewLowerDeck() *LowerDeck {
	return &LowerDeck{
		name:    "LowerDeck",
		options: make(map[string]interface{}),
	}
}

// SetName ...
func (ld *LowerDeck) SetName(name string) DataProcessor {
	ld.name = name
	return ld
}

// GetName ...
func (ld *LowerDeck) GetName() string {
	return ld.name
}

// SetOption ...
func (ld *LowerDeck) SetOption(name string, value interface{}) DataProcessor {
	ld.options[name] = value
	return ld
}

// GetOption ...
func (ld *LowerDeck) GetOption(name string) interface{} {
	if value, ok := ld.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable ...
func (ld *LowerDeck) SetEnable(enable bool) DataProcessor {
	return ld
}

// ForServer ...
func (ld *LowerDeck) ForServer(forServer bool) DataProcessor {
	return ld
}

// SetUStack ...
func (ld *LowerDeck) SetUStack(u UStack) DataProcessor {
	ld.ustack = u
	return ld
}

// SetUpperDataProcessor ...
func (ld *LowerDeck) SetUpperDataProcessor(dp DataProcessor) DataProcessor {
	ld.upperDataProcessor = dp
	return ld
}

// SetLowerDataProcessor ...
func (ld *LowerDeck) SetLowerDataProcessor(dp DataProcessor) DataProcessor {
	ld.lowerDataProcessor = dp
	return ld
}

// OnUpperPush ...
func (ld *LowerDeck) OnUpperPush(context Context) {
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

// OnLowerPush ...
func (ld *LowerDeck) OnLowerPush(context Context) {
}

// OnEvent ...
func (ld *LowerDeck) OnEvent(event Event) {
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

					ld.upperDataProcessor.OnLowerPush(
						NewUStackContext().
							SetConnection(connection).
							SetBuffer(ub))
				}
			}()
		}
	}()

	return ld
}
