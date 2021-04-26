// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// BytesCodec ...
type BytesCodec struct {
	ProcBase
}

// NewBytesCodec ...
func NewBytesCodec() DataProcessor {
	bc := &BytesCodec{
		NewProcBaseInstance("BytesCodec"),
	}
	return bc.ProcBase.SetWhere(bc)
}

// OnUpperData ...
func (bc *BytesCodec) OnUpperData(context Context) {
	if bc.enable {
		message := context.GetMessage()
		if message == nil {
			return
		}

		bytes, ok := message.([]byte)
		if !ok {
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
			return
		}

		bytes := make([]byte, ub.ReadableLength())

		n, err := ub.Read(bytes)
		if n == 0 || err != nil {
			return
		}

		context.SetMessage(bytes)
	}

	bc.upper.OnLowerData(context)
}
