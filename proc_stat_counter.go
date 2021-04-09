// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// StatCounter ...
type StatCounter struct {
	name               string
	enable             bool
	options            map[string]interface{}
	upperDataProcessor DataProcessor
	lowerDataProcessor DataProcessor
}

// NewStatCounter ...
func NewStatCounter() DataProcessor {
	return &StatCounter{
		name:    "StatCounter",
		enable:  true,
		options: make(map[string]interface{}),
	}
}

// SetName ...
func (sc *StatCounter) SetName(name string) DataProcessor {
	sc.name = name
	return sc
}

// GetName ...
func (sc *StatCounter) GetName() string {
	return sc.name
}

// SetOption ...
func (sc *StatCounter) SetOption(name string, value interface{}) DataProcessor {
	sc.options[name] = value
	return sc
}

// GetOption ...
func (sc *StatCounter) GetOption(name string) interface{} {
	if value, ok := sc.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable ...
func (sc *StatCounter) SetEnable(enable bool) DataProcessor {
	sc.enable = enable
	return sc
}

// ForServer ...
func (sc *StatCounter) ForServer(forServer bool) DataProcessor {
	return sc
}

// SetUStack ...
func (sc *StatCounter) SetUStack(u UStack) DataProcessor {
	return sc
}

// SetUpperDataProcessor ...
func (sc *StatCounter) SetUpperDataProcessor(dp DataProcessor) DataProcessor {
	sc.upperDataProcessor = dp
	return sc
}

// SetLowerDataProcessor ...
func (sc *StatCounter) SetLowerDataProcessor(dp DataProcessor) DataProcessor {
	sc.lowerDataProcessor = dp
	return sc
}

// OnUpperPush ...
func (sc *StatCounter) OnUpperPush(context Context) {
	if sc.enable {
		fmt.Println("StatCounter OnUpperPush")
	}

	sc.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (sc *StatCounter) OnLowerPush(context Context) {
	if sc.enable {
		fmt.Println("StatCounter OnLowerPush")
	}

	sc.upperDataProcessor.OnLowerPush(context)
}

// OnEvent ...
func (sc *StatCounter) OnEvent(event Event) {
}

// Run ...
func (sc *StatCounter) Run() DataProcessor {
	return sc
}
