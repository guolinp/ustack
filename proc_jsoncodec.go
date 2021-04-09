// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// JSONCodec ...
type JSONCodec struct {
	name               string
	enable             bool
	options            map[string]interface{}
	objectType         reflect.Type
	upperDataProcessor DataProcessor
	lowerDataProcessor DataProcessor
}

// NewJSONCodec ...
func NewJSONCodec(t reflect.Type) DataProcessor {
	return &JSONCodec{
		name:       "JSONCodec",
		enable:     true,
		options:    make(map[string]interface{}),
		objectType: t,
	}
}

// SetName ...
func (jc *JSONCodec) SetName(name string) DataProcessor {
	jc.name = name
	return jc
}

// GetName ...
func (jc *JSONCodec) GetName() string {
	return jc.name
}

// SetOption ...
func (jc *JSONCodec) SetOption(name string, value interface{}) DataProcessor {
	jc.options[name] = value
	return jc
}

// GetOption ...
func (jc *JSONCodec) GetOption(name string) interface{} {
	if value, ok := jc.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable ...
func (jc *JSONCodec) SetEnable(enable bool) DataProcessor {
	jc.enable = enable
	return jc
}

// ForServer ...
func (jc *JSONCodec) ForServer(forServer bool) DataProcessor {
	return jc
}

// SetUStack ...
func (jc *JSONCodec) SetUStack(u UStack) DataProcessor {
	return jc
}

// SetUpperDataProcessor ...
func (jc *JSONCodec) SetUpperDataProcessor(dp DataProcessor) DataProcessor {
	jc.upperDataProcessor = dp
	return jc
}

// SetLowerDataProcessor ...
func (jc *JSONCodec) SetLowerDataProcessor(dp DataProcessor) DataProcessor {
	jc.lowerDataProcessor = dp
	return jc
}

// OnUpperPush ...
func (jc *JSONCodec) OnUpperPush(context Context) {
	if jc.enable {
		message, ok := context.GetOption("message")

		if message == nil || !ok {
			fmt.Println("invalid message data")
			return
		}

		jsonBytes, err := json.Marshal(message)
		if err != nil {
			fmt.Println("failed to json marshal", err)
			return
		}

		// for test: 4096, 12
		ub := UBufAllocWithHeadReserved(4096, 128)

		n, err := ub.Write(jsonBytes)
		if n == 0 || err != nil {
			return
		}

		context.SetBuffer(ub)
	}

	jc.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (jc *JSONCodec) OnLowerPush(context Context) {
	if jc.enable {
		ub := context.GetBuffer()
		if ub == nil {
			fmt.Println("invalid lower data")
			return
		}

		size := ub.ReadableLength()
		data := make([]byte, size)
		n, err := ub.Read(data)

		if n == 0 || err != nil {
			return
		}

		objectItf := reflect.New(jc.objectType).Interface()
		err = json.Unmarshal(data, objectItf)
		if err != nil {
			fmt.Println("failed to json marshal", err)
			return
		}

		context.SetOption("message", objectItf)
	}

	jc.upperDataProcessor.OnLowerPush(context)
}

// OnEvent ...
func (jc *JSONCodec) OnEvent(event Event) {
}

// Run ...
func (jc *JSONCodec) Run() DataProcessor {
	return jc
}
