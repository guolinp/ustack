// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Loopback ...
type Loopback struct {
	Base
}

// NewLoopback returns a new instance
func NewLoopback() DataProcessor {
	lb := &Loopback{
		NewBaseInstance("Loopback"),
	}
	return lb.Base.SetWhere(lb)
}

// OnUpperPush sends back the message
func (lb *Loopback) OnUpperPush(context Context) {
	fmt.Println("Loopback: send back the uplayer data")
	lb.upperDataProcessor.OnLowerPush(context)
}
