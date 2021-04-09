// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// LoadBalancer ...
type LoadBalancer struct {
	name               string
	enable             bool
	options            map[string]interface{}
	upperDataProcessor DataProcessor
	lowerDataProcessor DataProcessor
}

// NewLoadBalancer ...
func NewLoadBalancer() DataProcessor {
	return &LoadBalancer{
		name:    "LoadBalancer",
		enable:  true,
		options: make(map[string]interface{}),
	}
}

// SetName ...
func (lb *LoadBalancer) SetName(name string) DataProcessor {
	lb.name = name
	return lb
}

// GetName ...
func (lb *LoadBalancer) GetName() string {
	return lb.name
}

// SetOption ...
func (lb *LoadBalancer) SetOption(name string, value interface{}) DataProcessor {
	lb.options[name] = value
	return lb
}

// GetOption ...
func (lb *LoadBalancer) GetOption(name string) interface{} {
	if value, ok := lb.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable ...
func (lb *LoadBalancer) SetEnable(enable bool) DataProcessor {
	lb.enable = enable
	return lb
}

// ForServer ...
func (lb *LoadBalancer) ForServer(forServer bool) DataProcessor {
	return lb
}

// SetUStack ...
func (lb *LoadBalancer) SetUStack(u UStack) DataProcessor {
	return lb
}

// SetUpperDataProcessor ...
func (lb *LoadBalancer) SetUpperDataProcessor(dp DataProcessor) DataProcessor {
	lb.upperDataProcessor = dp
	return lb
}

// SetLowerDataProcessor ...
func (lb *LoadBalancer) SetLowerDataProcessor(dp DataProcessor) DataProcessor {
	lb.lowerDataProcessor = dp
	return lb
}

// OnUpperPush ...
func (lb *LoadBalancer) OnUpperPush(context Context) {
	if lb.enable {
		fmt.Println("LoadBalancer OnUpperPush")
	}

	lb.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (lb *LoadBalancer) OnLowerPush(context Context) {
	if lb.enable {
		fmt.Println("LoadBalancer OnLowerPush")
	}

	lb.upperDataProcessor.OnLowerPush(context)
}

// OnEvent ...
func (lb *LoadBalancer) OnEvent(event Event) {
}

// Run ...
func (lb *LoadBalancer) Run() DataProcessor {
	return lb
}
