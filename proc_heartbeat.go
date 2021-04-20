// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import (
	"fmt"
	"sync"
	"time"
)

const (
	HeartbeatSelfMessageTag    byte = 0x83
	HeartbeatUplayerMessageTag byte = 0x00
)

// Heartbeat ...
type Heartbeat struct {
	ProcBase
	intervalInSecond int
	timeoutInSecond  int
	closeOnLost      bool
	mutex            sync.Mutex
	monitors         map[TransportConnection]time.Time
}

// NewHeartbeat ...
func NewHeartbeat() DataProcessor {
	hb := &Heartbeat{
		ProcBase:         NewProcBaseInstance("Heartbeat"),
		intervalInSecond: 1,
		timeoutInSecond:  30,
		closeOnLost:      true,
		monitors:         make(map[TransportConnection]time.Time, 16),
	}
	return hb.ProcBase.SetWhere(hb)
}

// updateMonitor ...
func (hb *Heartbeat) updateMonitor(connection TransportConnection) {
	hb.mutex.Lock()
	defer hb.mutex.Unlock()

	hb.monitors[connection] = time.Now()
}

// deleteMonitor ...
func (hb *Heartbeat) deleteMonitor(connection TransportConnection) {
	hb.mutex.Lock()
	defer hb.mutex.Unlock()

	_, ok := hb.monitors[connection]
	if ok {
		delete(hb.monitors, connection)
	}
}

// deleteMonitor ...
func (hb *Heartbeat) check() {
	hb.mutex.Lock()
	defer hb.mutex.Unlock()

	if len(hb.monitors) == 0 {
		return
	}

	for connection, lastTime := range hb.monitors {
		if int(time.Since(lastTime).Seconds()) < hb.timeoutInSecond {
			continue
		}

		fmt.Println("Heartbeat: connection", connection.GetName(), "lost")

		delete(hb.monitors, connection)

		if hb.closeOnLost {
			connection.Close()
		}

		hb.ustack.PublishEvent(Event{
			Type:   UStackEventHeartbeatLost,
			Source: hb,
			Data:   connection,
		})
	}
}

// GetOverhead returns the overhead
func (hb *Heartbeat) GetOverhead() int {
	return 1
}

// OnUpperData ...
func (hb *Heartbeat) OnUpperData(context Context) {
	ub := context.GetBuffer()
	if ub == nil {
		return
	}

	ub.WriteHeadByte(HeartbeatUplayerMessageTag)

	hb.lower.OnUpperData(context)
}

// OnLowerData ...
func (hb *Heartbeat) OnLowerData(context Context) {
	ub := context.GetBuffer()
	if ub == nil {
		return
	}

	tag, err := ub.ReadByte()
	if err != nil {
		return
	}

	if tag == HeartbeatSelfMessageTag {
		hb.updateMonitor(context.GetConnection())

		fmt.Printf("Heartbeat: %s, receive heartbeat\n", hb.GetName())
		return
	}

	hb.upper.OnLowerData(context)
}

// OnEvent ...
func (hb *Heartbeat) OnEvent(event Event) {
	if event.Type == UStackEventNewConnection {
		interval := hb.intervalInSecond
		connection := event.Data.(TransportConnection)

		hb.updateMonitor(connection)

		go func() {
			for {
				if connection.Closed() {
					fmt.Printf("Heartbeat: connection %s is closed\n", connection.GetName())
					hb.deleteMonitor(connection)
					return
				}

				ub := UBufAllocWithHeadReserved(
					hb.ustack.GetMTU(),
					hb.ustack.GetOverhead())

				ub.WriteByte(HeartbeatSelfMessageTag)

				hb.lower.OnUpperData(
					NewUStackContext().
						SetConnection(connection).
						SetBuffer(ub))

				fmt.Printf("Heartbeat: %s, send heartbeat\n", hb.GetName())

				time.Sleep(time.Second * time.Duration(interval))
			}
		}()
	}
}

// Run ...
func (hb *Heartbeat) Run() DataProcessor {
	interval, exists := OptionParseInt(hb.GetOption("IntervalInSecond"), hb.intervalInSecond)
	hb.intervalInSecond = interval
	if exists {
		fmt.Println("Heartbeat: option IntervalInSecond:", hb.intervalInSecond)
	}

	timeout, exists := OptionParseInt(hb.GetOption("TimeoutInSecond"), hb.timeoutInSecond)
	hb.timeoutInSecond = timeout
	if exists {
		fmt.Println("Heartbeat: option TimeoutInSecond:", hb.timeoutInSecond)
	}

	closeOnLost, exists := OptionParseBool(hb.GetOption("CloseOnLost"), hb.closeOnLost)
	hb.closeOnLost = closeOnLost
	if exists {
		fmt.Println("Heartbeat: option CloseOnLost:", hb.timeoutInSecond)
	}

	go func() {
		for {
			hb.check()
			time.Sleep(time.Second)
		}
	}()

	return hb
}
