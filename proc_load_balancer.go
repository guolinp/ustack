// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// LoadBalancer ...
type LoadBalancer struct {
	Base
}

// NewLoadBalancer ...
func NewLoadBalancer() DataProcessor {
	lb := &LoadBalancer{
		NewBaseInstance("LoadBalancer"),
	}
	return lb.Base.SetWhere(lb)
}

// OnUpperPush ...
func (lb *LoadBalancer) OnUpperPush(context Context) {
	if lb.enable {
		fmt.Println("LoadBalancer: OnUpperPush: todo")
	}

	lb.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (lb *LoadBalancer) OnLowerPush(context Context) {
	if lb.enable {
		fmt.Println("LoadBalancer: OnLowerPush: todo")
	}

	lb.upperDataProcessor.OnLowerPush(context)
}
