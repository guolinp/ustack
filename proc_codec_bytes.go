// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import (
	"fmt"
)

// BytesCodec ...
type BytesCodec struct {
	Base
}

// NewBytesCodec ...
func NewBytesCodec() DataProcessor {
	bc := &BytesCodec{
		NewBaseInstance("BytesCodec"),
	}
	return bc.Base.SetWhere(bc)
}

// OnUpperData ...
func (bc *BytesCodec) OnUpperData(context Context) {
	if bc.enable {
		bytes, ok := OptionParseByteSlice(context.GetOption("message"))
		if !ok {
			fmt.Println("BytesCodec: invalid uplayer message")
			return
		}

		ub := UBufAllocWithHeadReserved(
			bc.ustack.GetMTU(),
			bc.ustack.GetOverhead())

		n, err := ub.Write(bytes)
		if n == 0 || err != nil {
			return
		}

		context.SetBuffer(ub)
	}

	bc.lower.OnUpperData(context)
}

// OnLowerData ...
func (bc *BytesCodec) OnLowerData(context Context) {
	if bc.enable {
		ub := context.GetBuffer()
		if ub == nil {
			fmt.Println("BytesCodec: invalid lowlayer data")
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

	bc.upper.OnLowerData(context)
}
