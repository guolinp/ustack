// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// FrameDecoder ...
type FrameDecoder struct {
	name               string
	enable             bool
	options            map[string]interface{}
	upperDataProcessor DataProcessor
	lowerDataProcessor DataProcessor
}

// NewFrameDecoder ...
func NewFrameDecoder() DataProcessor {
	return &FrameDecoder{
		name:    "FrameDecoder",
		enable:  true,
		options: make(map[string]interface{}),
	}
}

// SetName ...
func (frm *FrameDecoder) SetName(name string) DataProcessor {
	frm.name = name
	return frm
}

// GetName ...
func (frm *FrameDecoder) GetName() string {
	return frm.name
}

// SetOption ...
func (frm *FrameDecoder) SetOption(name string, value interface{}) DataProcessor {
	frm.options[name] = value
	return frm
}

// GetOption ...
func (frm *FrameDecoder) GetOption(name string) interface{} {
	if value, ok := frm.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable ...
func (frm *FrameDecoder) SetEnable(enable bool) DataProcessor {
	frm.enable = enable
	return frm
}

// ForServer ...
func (frm *FrameDecoder) ForServer(forServer bool) DataProcessor {
	return frm
}

// SetUStack ...
func (frm *FrameDecoder) SetUStack(u UStack) DataProcessor {
	return frm
}

// SetUpperDataProcessor ...
func (frm *FrameDecoder) SetUpperDataProcessor(dp DataProcessor) DataProcessor {
	frm.upperDataProcessor = dp
	return frm
}

// SetLowerDataProcessor ...
func (frm *FrameDecoder) SetLowerDataProcessor(dp DataProcessor) DataProcessor {
	frm.lowerDataProcessor = dp
	return frm
}

// OnUpperPush ...
func (frm *FrameDecoder) OnUpperPush(context Context) {
	if frm.enable {
		fmt.Println("FrameDecoder OnUpperPush")
	}

	frm.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (frm *FrameDecoder) OnLowerPush(context Context) {
	if frm.enable {
		fmt.Println("FrameDecoder OnUpperPush")
	}

	frm.upperDataProcessor.OnLowerPush(context)
}

// OnEvent ...
func (frm *FrameDecoder) OnEvent(event Event) {
}

// Run ...
func (frm *FrameDecoder) Run() DataProcessor {
	return frm
}
