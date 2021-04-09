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
	name               string
	enable             bool
	options            map[string]interface{}
	objectType         reflect.Type
	upperDataProcessor DataProcessor
	lowerDataProcessor DataProcessor
}

// NewGOBCodec ...
func NewGOBCodec(t reflect.Type) DataProcessor {
	return &GOBCodec{
		name:       "GOBCodec",
		enable:     true,
		options:    make(map[string]interface{}),
		objectType: t,
	}
}

// SetName ...
func (g *GOBCodec) SetName(name string) DataProcessor {
	g.name = name
	return g
}

// GetName ...
func (g *GOBCodec) GetName() string {
	return g.name
}

// SetOption ...
func (g *GOBCodec) SetOption(name string, value interface{}) DataProcessor {
	g.options[name] = value
	return g
}

// GetOption ...
func (g *GOBCodec) GetOption(name string) interface{} {
	if value, ok := g.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable ...
func (g *GOBCodec) SetEnable(enable bool) DataProcessor {
	g.enable = enable
	return g
}

// ForServer ...
func (g *GOBCodec) ForServer(forServer bool) DataProcessor {
	return g
}

// SetUStack ...
func (g *GOBCodec) SetUStack(u UStack) DataProcessor {
	return g
}

// SetUpperDataProcessor ...
func (g *GOBCodec) SetUpperDataProcessor(dp DataProcessor) DataProcessor {
	g.upperDataProcessor = dp
	return g
}

// SetLowerDataProcessor ...
func (g *GOBCodec) SetLowerDataProcessor(dp DataProcessor) DataProcessor {
	g.lowerDataProcessor = dp
	return g
}

// OnUpperPush ...
func (g *GOBCodec) OnUpperPush(context Context) {
	if g.enable {
		message, ok := context.GetOption("message")

		if message == nil || !ok {
			fmt.Println("invalid message data")
			return
		}

		ub := UBufAllocWithHeadReserved(4096, 128)
		err := gob.NewEncoder(ub).Encode(message)
		if err != nil {
			fmt.Println("gob encode error:", err)
			return
		}

		context.SetBuffer(ub)
	}

	g.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (g *GOBCodec) OnLowerPush(context Context) {
	if g.enable {
		ub := context.GetBuffer()
		if ub == nil {
			fmt.Println("invalid lower data")
			return
		}

		objectItf := reflect.New(g.objectType).Interface()
		err := gob.NewDecoder(ub).Decode(objectItf)
		if err != nil {
			fmt.Println("gob encode error:", err)
			return
		}

		context.SetOption("message", objectItf)
	}

	g.upperDataProcessor.OnLowerPush(context)
}

// OnEvent ...
func (g *GOBCodec) OnEvent(event Event) {
}

// Run ...
func (g *GOBCodec) Run() DataProcessor {
	return g
}
