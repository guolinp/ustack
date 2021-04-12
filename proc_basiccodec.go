// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import (
	"fmt"
)

// BasicCodec ...
type BasicCodec struct {
	Base
}

// NewBasicCodec ...
func NewBasicCodec() DataProcessor {
	bc := &BasicCodec{
		NewBaseInstance("BasicCodec"),
	}
	return bc.Base.SetWhere(bc)
}

// OnUpperPush ...
func (bc *BasicCodec) OnUpperPush(context Context) {
	if bc.enable {
		message, ok := context.GetOption("message")

		if message == nil || !ok {
			fmt.Println("invalid message data")
			return
		}

		bytes, ok := message.([]byte)
		if !ok {
			return
		}

		// for test: 4096, 12
		ub := UBufAllocWithHeadReserved(4096, 128)

		n, err := ub.Write(bytes)
		if n == 0 || err != nil {
			return
		}

		context.SetBuffer(ub)
	}

	bc.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (bc *BasicCodec) OnLowerPush(context Context) {
	if bc.enable {
		ub := context.GetBuffer()
		if ub == nil {
			fmt.Println("invalid lower data")
			return
		}

		size := ub.ReadableLength()
		bytes := make([]byte, size)
		n, err := ub.Read(bytes)

		if n == 0 || err != nil {
			return
		}

		context.SetOption("message", bytes)
	}

	bc.upperDataProcessor.OnLowerPush(context)
}
