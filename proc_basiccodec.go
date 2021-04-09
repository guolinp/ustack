// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import (
	"fmt"
)

// BasicCodec ...
type BasicCodec struct {
	name               string
	enable             bool
	options            map[string]interface{}
	upperDataProcessor DataProcessor
	lowerDataProcessor DataProcessor
}

// NewBasicCodec ...
func NewBasicCodec() DataProcessor {
	return &BasicCodec{
		name:    "BasicCodec",
		enable:  true,
		options: make(map[string]interface{}),
	}
}

// SetName ...
func (bc *BasicCodec) SetName(name string) DataProcessor {
	bc.name = name
	return bc
}

// GetName ...
func (bc *BasicCodec) GetName() string {
	return bc.name
}

// SetOption ...
func (bc *BasicCodec) SetOption(name string, value interface{}) DataProcessor {
	bc.options[name] = value
	return bc
}

// GetOption ...
func (bc *BasicCodec) GetOption(name string) interface{} {
	if value, ok := bc.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable ...
func (bc *BasicCodec) SetEnable(enable bool) DataProcessor {
	bc.enable = enable
	return bc
}

// ForServer ...
func (bc *BasicCodec) ForServer(forServer bool) DataProcessor {
	return bc
}

// SetUStack ...
func (bc *BasicCodec) SetUStack(u UStack) DataProcessor {
	return bc
}

// SetUpperDataProcessor ...
func (bc *BasicCodec) SetUpperDataProcessor(dp DataProcessor) DataProcessor {
	bc.upperDataProcessor = dp
	return bc
}

// SetLowerDataProcessor ...
func (bc *BasicCodec) SetLowerDataProcessor(dp DataProcessor) DataProcessor {
	bc.lowerDataProcessor = dp
	return bc
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

// OnEvent ...
func (bc *BasicCodec) OnEvent(event Event) {
}

// Run ...
func (bc *BasicCodec) Run() DataProcessor {
	return bc
}
