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

// OnUpperData ...
func (jc *JSONCodec) OnUpperData(context Context) {
	if jc.enable {
		message := context.GetOption("message")
		if message == nil {
			return
		}

		jsonBytes, err := json.Marshal(message)
		if err != nil {
			fmt.Println("JSONCodec: failed to json marshal", err)
			return
		}

		ub := UBufAllocWithHeadReserved(
			jc.ustack.GetMTU(),
			jc.ustack.GetOverhead())

		n, err := ub.Write(jsonBytes)
		if n == 0 || err != nil {
			return
		}

		context.SetBuffer(ub)
	}

	jc.lower.OnUpperData(context)
}

// OnLowerData ...
func (jc *JSONCodec) OnLowerData(context Context) {
	if jc.enable {
		ub := context.GetBuffer()
		if ub == nil {
			return
		}

		data := make([]byte, ub.ReadableLength())
		
		n, err := ub.Read(data)
		if n == 0 || err != nil {
			return
		}

		objectItf := reflect.New(jc.objectType).Interface()
		err = json.Unmarshal(data, objectItf)
		if err != nil {
			fmt.Println("JSONCodec: failed to json marshal", err)
			return
		}

		context.SetOption("message", objectItf)
	}

	jc.upper.OnLowerData(context)
}
