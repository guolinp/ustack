// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// StringCodec ...
type StringCodec struct {
	ProcBase
}

// NewStringCodec ...
func NewStringCodec() DataProcessor {
	bc := &StringCodec{
		NewProcBaseInstance("StringCodec"),
	}
	return bc.ProcBase.SetWhere(bc)
}

// OnUpperData ...
func (bc *StringCodec) OnUpperData(context Context) {

	if bc.enable {
		message := context.GetMessage()
		if message == nil {
			return
		}

		str, ok := message.(string)
		if !ok {
			return
		}

		ub := UBufAllocWithHeadReserved(
			bc.ustack.GetMTU(),
			bc.ustack.GetOverhead())

		n, err := ub.Write([]byte(str))
		if n == 0 || err != nil {
			return
		}

		context.SetBuffer(ub)
	}

	bc.lower.OnUpperData(context)
}

// OnLowerData ...
func (bc *StringCodec) OnLowerData(context Context) {
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

		context.SetMessage(string(bytes))
	}

	bc.upper.OnLowerData(context)
}
