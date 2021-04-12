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
	Base
	objectType reflect.Type
}

// NewGOBCodec ...
func NewGOBCodec(t reflect.Type) DataProcessor {
	g := &GOBCodec{
		Base:       NewBaseInstance("GOBCodec"),
		objectType: t,
	}
	return g.Base.SetWhere(g)
}

// OnUpperData ...
func (g *GOBCodec) OnUpperData(context Context) {
	if g.enable {
		message, ok := context.GetOption("message")

		if message == nil || !ok {
			fmt.Println("GOBCodec: invalid message data")
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
			fmt.Println("GOBCodec: invalid lower data")
			return
		}

		objectItf := reflect.New(g.objectType).Interface()
		err := gob.NewDecoder(ub).Decode(objectItf)
		if err != nil {
			fmt.Println("GOBCodec: gob encode error:", err)
			return
		}

		context.SetOption("message", objectItf)
	}

	g.upper.OnLowerData(context)
}
