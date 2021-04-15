// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// Feature ...
type Feature interface {
	SetName(name string) Feature
	GetName() string
	SetOption(name string, value interface{}) Feature
	GetOption(name string) interface{}
	SetUStack(ustack UStack) Feature
	OnEvent(event Event)
	Run() Feature
}
