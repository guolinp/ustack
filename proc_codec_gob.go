// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import (
	"encoding/gob"
	"fmt"
	"reflect"
)

// GOBCodec ...
type GOBCodec struct {
	ProcBase
	objectType reflect.Type
}

// NewGOBCodec ...
func NewGOBCodec(t reflect.Type) DataProcessor {
	g := &GOBCodec{
		ProcBase:       NewProcBaseInstance("GOBCodec"),
		objectType: t,
	}
	return g.ProcBase.SetWhere(g)
}

// OnUpperData ...
func (g *GOBCodec) OnUpperData(context Context) {
	if g.enable {
		message := context.GetMessage()
		if message == nil {
			return
		}

		ub := UBufAllocWithHeadReserved(
			g.ustack.GetMTU(),
			g.ustack.GetOverhead())

		err := gob.NewEncoder(ub).Encode(message)
		if err != nil {
			fmt.Println("GOBCodec: gob encode error:", err)
			return
		}

		context.SetBuffer(ub)
	}

	g.lower.OnUpperData(context)
}

// OnLowerData ...
func (g *GOBCodec) OnLowerData(context Context) {
	if g.enable {
		ub := context.GetBuffer()
		if ub == nil {
			return
		}

		objectItf := reflect.New(g.objectType).Interface()
		err := gob.NewDecoder(ub).Decode(objectItf)
		if err != nil {
			fmt.Println("GOBCodec: gob encode error:", err)
			return
		}

		context.SetMessage(objectItf)
	}

	g.upper.OnLowerData(context)
}
