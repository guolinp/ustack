// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Loopback ...
type Loopback struct {
	ProcBase
}

// NewLoopback returns a new instance
func NewLoopback() DataProcessor {
	lb := &Loopback{
		NewProcBaseInstance("Loopback"),
	}
	return lb.ProcBase.SetWhere(lb)
}

// OnUpperData sends back the message
func (lb *Loopback) OnUpperData(context Context) {
	fmt.Println("Loopback: send back the uplayer data")
	lb.upper.OnLowerData(context)
}
