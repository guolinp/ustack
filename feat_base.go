// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

type FeatBase struct {
	where   Feature
	name    string
	ustack  UStack
	options map[string]interface{}
}

// NewFeatBaseInstance returns a new instance
func NewFeatBaseInstance(name string) FeatBase {
	base := FeatBase{
		name:    name,
		ustack:  nil,
		options: make(map[string]interface{}),
	}
	// by default is itself
	base.where = &base
	return base
}

// NewFeatBase return a new instance that meets for Feature interface
func NewFeatBase() Feature {
	base := NewFeatBaseInstance("FeatBase")
	return &base
}

// SetName set the name
func (fb *FeatBase) SetName(name string) Feature {
	fb.name = name
	return fb.where
}

// GetName returns name
func (fb *FeatBase) GetName() string {
	return fb.name
}

// SetOption
func (fb *FeatBase) SetOption(name string, value interface{}) Feature {
	fb.options[name] = value
	return fb.where
}

// GetOption
func (fb *FeatBase) GetOption(name string) interface{} {
	if value, ok := fb.options[name]; ok {
		return value
	}
	return nil
}

// SetUStack set the UStack instance
func (fb *FeatBase) SetUStack(ustack UStack) Feature {
	fb.ustack = ustack
	return fb.where
}

// OnEvent is called when any event hanppen
func (base *FeatBase) OnEvent(event Event) {
}

// Run
func (fb *FeatBase) Run() Feature {
	return fb.where
}
