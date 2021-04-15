// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// LoadBalancer ...
type LoadBalancer struct {
	ProcBase
}

// NewLoadBalancer ...
func NewLoadBalancer() DataProcessor {
	lb := &LoadBalancer{
		NewProcBaseInstance("LoadBalancer"),
	}
	return lb.ProcBase.SetWhere(lb)
}

// OnUpperData ...
func (lb *LoadBalancer) OnUpperData(context Context) {
	if lb.enable {
		fmt.Println("LoadBalancer: OnUpperData: todo")
	}

	lb.lower.OnUpperData(context)
}

// OnLowerData ...
func (lb *LoadBalancer) OnLowerData(context Context) {
	if lb.enable {
		fmt.Println("LoadBalancer: OnLowerData: todo")
	}

	lb.upper.OnLowerData(context)
}
