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
	Base
	objectType reflect.Type
}

// NewJSONCodec ...
func NewJSONCodec(t reflect.Type) DataProcessor {
	jc := &JSONCodec{
		Base:       NewBaseInstance("JSONCodec"),
		objectType: t,
	}
	return jc.Base.SetWhere(jc)
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
