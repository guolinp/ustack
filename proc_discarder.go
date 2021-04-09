// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Discarder ...
type Discarder struct {
	name               string
	enable             bool
	options            map[string]interface{}
	upperDataProcessor DataProcessor
	lowerDataProcessor DataProcessor
}

// NewDiscarder ...
func NewDiscarder() DataProcessor {
	return &Discarder{
		name:    "Discarder",
		enable:  true,
		options: make(map[string]interface{}),
	}
}

// SetName ...
func (dis *Discarder) SetName(name string) DataProcessor {
	dis.name = name
	return dis
}

// GetName ...
func (dis *Discarder) GetName() string {
	return dis.name
}

// SetOption ...
func (dis *Discarder) SetOption(name string, value interface{}) DataProcessor {
	dis.options[name] = value
	return dis
}

// GetOption ...
func (dis *Discarder) GetOption(name string) interface{} {
	if value, ok := dis.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable ...
func (dis *Discarder) SetEnable(enable bool) DataProcessor {
	dis.enable = enable
	return dis
}

// ForServer ...
func (dis *Discarder) ForServer(forServer bool) DataProcessor {
	return dis
}

// SetUStack ...
func (dis *Discarder) SetUStack(u UStack) DataProcessor {
	return dis
}

// SetUpperDataProcessor ...
func (dis *Discarder) SetUpperDataProcessor(dp DataProcessor) DataProcessor {
	dis.upperDataProcessor = dp
	return dis
}

// SetLowerDataProcessor ...
func (dis *Discarder) SetLowerDataProcessor(dp DataProcessor) DataProcessor {
	dis.lowerDataProcessor = dp
	return dis
}

// OnUpperPush ...
func (dis *Discarder) OnUpperPush(context Context) {
	if dis.enable {
		fmt.Println("Discarder OnUpperPush")
	}

	dis.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (dis *Discarder) OnLowerPush(context Context) {
	if dis.enable {
		fmt.Println("Discarder OnLowerPush")
	}

	dis.upperDataProcessor.OnLowerPush(context)
}

// OnEvent ...
func (dis *Discarder) OnEvent(event Event) {
}

// Run ...
func (dis *Discarder) Run() DataProcessor {
	return dis
}
