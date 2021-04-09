// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Loopback ...
type Loopback struct {
	name               string
	enable             bool
	options            map[string]interface{}
	upperDataProcessor DataProcessor
	lowerDataProcessor DataProcessor
}

// NewLoopback ...
func NewLoopback() DataProcessor {
	return &Loopback{
		name:    "Loopback",
		enable:  true,
		options: make(map[string]interface{}),
	}
}

// SetName ...
func (lb *Loopback) SetName(name string) DataProcessor {
	lb.name = name
	return lb
}

// GetName ...
func (lb *Loopback) GetName() string {
	return lb.name
}

// SetOption ...
func (lb *Loopback) SetOption(name string, value interface{}) DataProcessor {
	lb.options[name] = value
	return lb
}

// GetOption ...
func (lb *Loopback) GetOption(name string) interface{} {
	if value, ok := lb.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable ...
func (lb *Loopback) SetEnable(enable bool) DataProcessor {
	lb.enable = enable
	return lb
}

// ForServer ...
func (lb *Loopback) ForServer(forServer bool) DataProcessor {
	return lb
}

// SetUStack ...
func (lb *Loopback) SetUStack(u UStack) DataProcessor {
	return lb
}

// SetUpperDataProcessor ...
func (lb *Loopback) SetUpperDataProcessor(dp DataProcessor) DataProcessor {
	lb.upperDataProcessor = dp
	return lb
}

// SetLowerDataProcessor ...
func (lb *Loopback) SetLowerDataProcessor(dp DataProcessor) DataProcessor {
	lb.lowerDataProcessor = dp
	return lb
}

// OnUpperPush ...
func (lb *Loopback) OnUpperPush(context Context) {
	if lb.enable {
		fmt.Println("Loopback OnUpperPush")
	}

	lb.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (lb *Loopback) OnLowerPush(context Context) {
	if lb.enable {
		fmt.Println("Loopback OnLowerPush")
	}

	lb.upperDataProcessor.OnLowerPush(context)
}

// OnEvent ...
func (lb *Loopback) OnEvent(event Event) {
}

// Run ...
func (lb *Loopback) Run() DataProcessor {
	return lb
}
