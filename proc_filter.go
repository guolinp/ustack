// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Filter ...
type Filter struct {
	name               string
	enable             bool
	options            map[string]interface{}
	upperDataProcessor DataProcessor
	lowerDataProcessor DataProcessor
}

// NewFilter ...
func NewFilter() DataProcessor {
	return &Filter{
		name:    "Filter",
		enable:  true,
		options: make(map[string]interface{}),
	}
}

// SetName ...
func (filter *Filter) SetName(name string) DataProcessor {
	filter.name = name
	return filter
}

// GetName ...
func (filter *Filter) GetName() string {
	return filter.name
}

// SetOption ...
func (filter *Filter) SetOption(name string, value interface{}) DataProcessor {
	filter.options[name] = value
	return filter
}

// GetOption ...
func (filter *Filter) GetOption(name string) interface{} {
	if value, ok := filter.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable ...
func (filter *Filter) SetEnable(enable bool) DataProcessor {
	filter.enable = enable
	return filter
}

// ForServer ...
func (filter *Filter) ForServer(forServer bool) DataProcessor {
	return filter
}

// SetUStack ...
func (filter *Filter) SetUStack(u UStack) DataProcessor {
	return filter
}

// SetUpperDataProcessor ...
func (filter *Filter) SetUpperDataProcessor(dp DataProcessor) DataProcessor {
	filter.upperDataProcessor = dp
	return filter
}

// SetLowerDataProcessor ...
func (filter *Filter) SetLowerDataProcessor(dp DataProcessor) DataProcessor {
	filter.lowerDataProcessor = dp
	return filter
}

// OnUpperPush ...
func (filter *Filter) OnUpperPush(context Context) {
	if filter.enable {
		fmt.Println("Filter OnUpperPush")
	}

	filter.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (filter *Filter) OnLowerPush(context Context) {
	if filter.enable {
		fmt.Println("Filter OnLowerPush")
	}

	filter.upperDataProcessor.OnLowerPush(context)
}

// OnEvent ...
func (filter *Filter) OnEvent(event Event) {
}

// Run ...
func (filter *Filter) Run() DataProcessor {
	return filter
}
