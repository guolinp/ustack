// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Echo ...
type Echo struct {
	name               string
	enable             bool
	options            map[string]interface{}
	upperDataProcessor DataProcessor
	lowerDataProcessor DataProcessor
}

// NewEcho ...
func NewEcho() DataProcessor {
	return &Echo{
		name:    "Echo",
		enable:  true,
		options: make(map[string]interface{}),
	}
}

// SetName ...
func (echo *Echo) SetName(name string) DataProcessor {
	echo.name = name
	return echo
}

// GetName ...
func (echo *Echo) GetName() string {
	return echo.name
}

// SetOption ...
func (echo *Echo) SetOption(name string, value interface{}) DataProcessor {
	echo.options[name] = value
	return echo
}

// GetOption ...
func (echo *Echo) GetOption(name string) interface{} {
	if value, ok := echo.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable ...
func (echo *Echo) SetEnable(enable bool) DataProcessor {
	echo.enable = enable
	return echo
}

// ForServer ...
func (echo *Echo) ForServer(forServer bool) DataProcessor {
	return echo
}

// SetUStack ...
func (echo *Echo) SetUStack(u UStack) DataProcessor {
	return echo
}

// SetUpperDataProcessor ...
func (echo *Echo) SetUpperDataProcessor(dp DataProcessor) DataProcessor {
	echo.upperDataProcessor = dp
	return echo
}

// SetLowerDataProcessor ...
func (echo *Echo) SetLowerDataProcessor(dp DataProcessor) DataProcessor {
	echo.lowerDataProcessor = dp
	return echo
}

// OnUpperPush ...
func (echo *Echo) OnUpperPush(context Context) {
	if echo.enable {
		fmt.Println("Echo OnUpperPush")
	}

	echo.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (echo *Echo) OnLowerPush(context Context) {
	if echo.enable {
		fmt.Println("Echo OnLowerPush")
	}

	echo.upperDataProcessor.OnLowerPush(context)
}

// OnEvent ...
func (echo *Echo) OnEvent(event Event) {
}

// Run ...
func (echo *Echo) Run() DataProcessor {
	return echo
}
