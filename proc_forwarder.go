// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Forwarder ...
type Forwarder struct {
	name               string
	enable             bool
	options            map[string]interface{}
	upperDataProcessor DataProcessor
	lowerDataProcessor DataProcessor
}

// NewForwarder ...
func NewForwarder() DataProcessor {
	return &Forwarder{
		name:    "Forwarder",
		enable:  true,
		options: make(map[string]interface{}),
	}
}

// SetName ...
func (fwd *Forwarder) SetName(name string) DataProcessor {
	fwd.name = name
	return fwd
}

// GetName ...
func (fwd *Forwarder) GetName() string {
	return fwd.name
}

// SetOption ...
func (fwd *Forwarder) SetOption(name string, value interface{}) DataProcessor {
	fwd.options[name] = value
	return fwd
}

// GetOption ...
func (fwd *Forwarder) GetOption(name string) interface{} {
	if value, ok := fwd.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable ...
func (fwd *Forwarder) SetEnable(enable bool) DataProcessor {
	fwd.enable = enable
	return fwd
}

// ForServer ...
func (fwd *Forwarder) ForServer(forServer bool) DataProcessor {
	return fwd
}

// SetUStack ...
func (fwd *Forwarder) SetUStack(u UStack) DataProcessor {
	return fwd
}

// SetUpperDataProcessor ...
func (fwd *Forwarder) SetUpperDataProcessor(dp DataProcessor) DataProcessor {
	fwd.upperDataProcessor = dp
	return fwd
}

// SetLowerDataProcessor ...
func (fwd *Forwarder) SetLowerDataProcessor(dp DataProcessor) DataProcessor {
	fwd.lowerDataProcessor = dp
	return fwd
}

// OnUpperPush ...
func (fwd *Forwarder) OnUpperPush(context Context) {
	if fwd.enable {
		fmt.Println("Forwarder OnUpperPush")
	}

	fwd.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (fwd *Forwarder) OnLowerPush(context Context) {
	if fwd.enable {
		fmt.Println("Forwarder OnLowerPush")
	}

	fwd.upperDataProcessor.OnLowerPush(context)
}

// OnEvent ...
func (fwd *Forwarder) OnEvent(event Event) {
}

// Run ...
func (fwd *Forwarder) Run() DataProcessor {
	return fwd
}
