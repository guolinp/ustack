// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import (
	"fmt"
)

// EncoderFn does encode message to UBuf
type EncoderFn func(message interface{}, ub *UBuf) error

// EncoderFn does decode from UBuf to message
type DecoderFn func(ub *UBuf) (interface{}, error)

// GenericCodec ...
type GenericCodec struct {
	Base
	encoder EncoderFn
	decoder DecoderFn
}

// NewGenericCodec ...
func NewGenericCodec(encoder EncoderFn, decoder DecoderFn) DataProcessor {
	gc := &GenericCodec{
		Base:    NewBaseInstance("GenericCodec"),
		encoder: encoder,
		decoder: decoder,
	}
	return gc.Base.SetWhere(gc)
}

// OnUpperData ...
func (gc *GenericCodec) OnUpperData(context Context) {
	if gc.enable {
		if gc.encoder == nil {
			fmt.Println("GenericCodec: not found the encoder")
			return
		}

		message, ok := context.GetOption("message")

		if message == nil || !ok {
			fmt.Println("GenericCodec: invalid message data")
			return
		}

		ub := UBufAllocWithHeadReserved(
			gc.ustack.GetMTU(),
			gc.ustack.GetOverhead())

		err := gc.encoder(message, ub)
		if err != nil {
			fmt.Println("GenericCodec: encode error:")
			return
		}

		context.SetBuffer(ub)
	}

	gc.lower.OnUpperData(context)
}

// OnLowerData ...
func (gc *GenericCodec) OnLowerData(context Context) {
	if gc.enable {
		if gc.decoder == nil {
			fmt.Println("GenericCodec: not found the decoder")
			return
		}

		ub := context.GetBuffer()
		if ub == nil {
			fmt.Println("GenericCodec: invalid lower data")
			return
		}

		message, err := gc.decoder(ub)
		if err != nil {
			fmt.Println("GenericCodec: encode error:", err)
			return
		}

		context.SetOption("message", message)
	}

	gc.upper.OnLowerData(context)
}