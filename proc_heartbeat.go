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
	Base
}

// NewHeartbeat ...
func NewHeartbeat() DataProcessor {
	hb := &Heartbeat{
		NewBaseInstance("Heartbeat"),
	}
	return hb.Base.SetWhere(hb)
}

// GetOverhead returns the overhead
func (hb *Heartbeat) GetOverhead() int {
	return 1
}

// OnUpperData ...
func (hb *Heartbeat) OnUpperData(context Context) {
	hb.lower.OnUpperData(context)
}

// OnLowerData ...
func (hb *Heartbeat) OnLowerData(context Context) {
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
	hb.upper.OnLowerData(context)
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

				hb.lower.OnUpperData(
					NewUStackContext().
						SetConnection(connection).
						SetBuffer(ub))

				fmt.Printf("Heartbeat: %s, send heartbeat\n", hb.GetName())
				time.Sleep(time.Second)
			}
		}()
	}
}
