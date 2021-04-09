// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import (
	"fmt"
	"time"
)

// Heartbeat ...
type Heartbeat struct {
	name               string
	enbale             bool
	options            map[string]interface{}
	upperDataProcessor DataProcessor
	lowerDataProcessor DataProcessor
}

// NewHeartbeat ...
func NewHeartbeat() DataProcessor {
	return &Heartbeat{
		name:    "Heartbeat",
		enbale:  true,
		options: make(map[string]interface{}),
	}
}

// SetName ...
func (hb *Heartbeat) SetName(name string) DataProcessor {
	hb.name = name
	return hb
}

// GetName ...
func (hb *Heartbeat) GetName() string {
	return hb.name
}

// SetOption ...
func (hb *Heartbeat) SetOption(name string, value interface{}) DataProcessor {
	hb.options[name] = value
	return hb
}

// GetOption ...
func (hb *Heartbeat) GetOption(name string) interface{} {
	if value, ok := hb.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable ...
func (hb *Heartbeat) SetEnable(enable bool) DataProcessor {
	hb.enbale = enable
	return hb
}

// ForServer ...
func (hb *Heartbeat) ForServer(forServer bool) DataProcessor {
	return hb
}

// SetUStack ...
func (hb *Heartbeat) SetUStack(u UStack) DataProcessor {
	return hb
}

// SetUpperDataProcessor ...
func (hb *Heartbeat) SetUpperDataProcessor(dp DataProcessor) DataProcessor {
	hb.upperDataProcessor = dp
	return hb
}

// SetLowerDataProcessor ...
func (hb *Heartbeat) SetLowerDataProcessor(dp DataProcessor) DataProcessor {
	hb.lowerDataProcessor = dp
	return hb
}

// OnUpperPush ...
func (hb *Heartbeat) OnUpperPush(context Context) {
	hb.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (hb *Heartbeat) OnLowerPush(context Context) {
	ub := context.GetBuffer()
	if ub == nil {
		return
	}

	heartbeat, err := ub.PeekByte()
	if err != nil {
		return
	}

	if heartbeat == 0x12 {
		fmt.Printf("Heartbeat: %s, receive heartbeat\n", hb.GetName())
		return
	}

	// not a heartbeat message, pass it to uplayer
	hb.upperDataProcessor.OnLowerPush(context)
}

// OnEvent ...
func (hb *Heartbeat) OnEvent(event Event) {
	if event.Type == UStackEventNewConnection {
		connection := event.Data.(TransportConnection)
		go func() {
			for {
				if connection.Closed() {
					fmt.Printf("Heartbeat: %s, connection is closed\n", hb.GetName())
					return
				}

				ub := UBufAlloc(1)
				ub.WriteByte(0x12)

				hb.lowerDataProcessor.OnUpperPush(NewUStackContext().
					SetConnection(connection).
					SetBuffer(ub))

				fmt.Printf("Heartbeat: %s, send heartbeat\n", hb.GetName())
				time.Sleep(time.Second)
			}
		}()
	}
}

// Run ...
func (hb *Heartbeat) Run() DataProcessor {
	return hb
}
