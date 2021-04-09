// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// DataProcessor ...
type DataProcessor interface {
	SetName(name string) DataProcessor
	GetName() string

	ForServer(bool) DataProcessor

	SetOption(name string, value interface{}) DataProcessor
	GetOption(name string) interface{}

	SetEnable(enable bool) DataProcessor

	SetUStack(u UStack) DataProcessor

	SetUpperDataProcessor(dp DataProcessor) DataProcessor
	SetLowerDataProcessor(dp DataProcessor) DataProcessor

	OnUpperPush(context Context)
	OnLowerPush(context Context)

	OnEvent(event Event)

	Run() DataProcessor
}
