// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import (
	"fmt"
	"time"
)

const (
	StatCounterSelfMessageReqTag byte = 0x10
	StatCounterSelfMessageResTag byte = 0x23
	StatCounterUplayerMessageTag byte = 0x00
)

// StatCounter ...
type StatCounter struct {
	ProcBase
	intervalInSecond int
	txCounter        uint64
	rxCounter        uint64
}

// NewStatCounter ...
func NewStatCounter() DataProcessor {
	sc := &StatCounter{
		ProcBase:         NewProcBaseInstance("StatCounter"),
		intervalInSecond: 0,
		txCounter:        0,
		rxCounter:        0,
	}
	return sc.ProcBase.SetWhere(sc)
}

// GetOverhead returns the overhead
func (hb *StatCounter) GetOverhead() int {
	return 1
}

// OnUpperData ...
func (sc *StatCounter) OnUpperData(context Context) {
	if context.GetConnection().UseReference() {
		sc.lower.OnUpperData(context)
		return
	}

	ub := context.GetBuffer()
	if ub == nil {
		return
	}

	if sc.enable {
		sc.txCounter++
		// fmt.Println("StatCounter: txCounter:", sc.txCounter)
	}

	ub.WriteHeadByte(StatCounterUplayerMessageTag)

	sc.lower.OnUpperData(context)
}

// request, just for test
func (sc *StatCounter) request(connection TransportConnection) {
	ub := UBufAllocWithHeadReserved(
		sc.ustack.GetMTU(),
		sc.ustack.GetOverhead())

	ub.WriteByte(StatCounterSelfMessageReqTag)

	sc.lower.OnUpperData(
		NewUStackContext().
			SetConnection(connection).
			SetBuffer(ub))
}

// response, just for test
func (sc *StatCounter) response(context Context) {
	ub := UBufAllocWithHeadReserved(
		sc.ustack.GetMTU(),
		sc.ustack.GetOverhead())

	ub.WriteHeadByte(StatCounterSelfMessageResTag)

	ub.WriteU64BE(sc.txCounter)
	ub.WriteU64BE(sc.rxCounter)

	ub.Write([]byte(sc.ustack.GetName()))

	context.SetBuffer(ub)
	sc.lower.OnUpperData(context)
}

// show, just for test
func (sc *StatCounter) show(context Context) {
	ub := context.GetBuffer()

	txCounter, _ := ub.ReadU64BE()
	rxCounter, _ := ub.ReadU64BE()

	name := make([]byte, 128)
	ub.Read(name)

	fmt.Println("\n######## UStack:", string(name), "txCounter:", txCounter, "rxCounter", rxCounter)
}

// OnLowerData ...
func (sc *StatCounter) OnLowerData(context Context) {
	if context.GetConnection().UseReference() {
		sc.upper.OnLowerData(context)
		return
	}

	ub := context.GetBuffer()
	if ub == nil {
		return
	}

	tag, err := ub.ReadByte()
	if err != nil {
		return
	}

	if sc.enable {
		if tag == StatCounterSelfMessageReqTag {
			// fmt.Printf("StatCounter: collect request received on connection: %s\n",
			// 	context.GetConnection().GetName())
			sc.response(context)
			return
		} else if tag == StatCounterSelfMessageResTag {
			// fmt.Printf("StatCounter: collect response received on connection: %s\n",
			// 	context.GetConnection().GetName())
			sc.show(context)
			return
		} else {
			sc.rxCounter++
			// fmt.Println("StatCounter: rxCounter:", sc.rxCounter)
		}
	}

	sc.upper.OnLowerData(context)
}

// OnEvent ...
func (sc *StatCounter) OnEvent(event Event) {
	if event.Type == UStackEventNewConnection {
		interval := sc.intervalInSecond

		if interval <= 0 {
			return
		}

		connection := event.Data.(TransportConnection)

		go func() {
			for {
				if connection.Closed() {
					fmt.Printf("StatCounter: connection %s is closed\n", connection.GetName())
					return
				}
				sc.request(connection)
				// fmt.Printf("StatCounter: send collect request on connection: %s\n",
				// 	connection.GetName())
				time.Sleep(time.Second * time.Duration(interval))
			}
		}()
	}
}

// Run ...
func (sc *StatCounter) Run() DataProcessor {
	interval, exists := OptionParseInt(sc.GetOption("Collect.IntervalInSecond"), sc.intervalInSecond)
	sc.intervalInSecond = interval
	if exists {
		fmt.Println("StatCounter: option Collect.IntervalInSecond:", sc.intervalInSecond)
	}
	return sc
}
