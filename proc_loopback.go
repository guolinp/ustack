// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Loopback ...
type Loopback struct {
	Base
}

// NewLoopback ...
func NewLoopback() DataProcessor {
	lb := &Loopback{
		NewBaseInstance("Loopback"),
	}
	return lb.Base.SetWhere(lb)
}

// OnUpperPush ...
func (lb *Loopback) OnUpperPush(context Context) {
	if lb.enable {
		fmt.Println("Loopback OnUpperPush")
	}

	lb.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (lb *Loopback) OnLowerPush(context Context) {
	if lb.enable {
		fmt.Println("Loopback OnLowerPush")
	}

	lb.upperDataProcessor.OnLowerPush(context)
}
