// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

const LoadBalancerFieldSizeInByte int = 1

const (
	LoadBalancerSelfMessageReqLoadTag byte = 0x19
	LoadBalancerSelfMessageResLoadTag byte = 0x82
	LoadBalancerSelfMessageReqWorkTag byte = 0x11
	LoadBalancerSelfMessageResWorkTag byte = 0x08
	LoadBalancerUplayerMessageTag     byte = 0x00
)

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

// GetOverhead ...
func (lb *LoadBalancer) GetOverhead() int {
	// for saving load balancer field
	return LoadBalancerFieldSizeInByte
}

// OnUpperData ...
func (lb *LoadBalancer) OnUpperData(context Context) {
	ub := context.GetBuffer()
	if ub == nil {
		return
	}

	if lb.enable {
		fmt.Println("LoadBalancer: OnUpperData: todo")
		ub.WriteHeadByte(LoadBalancerSelfMessageReqWorkTag)
	} else {
		ub.WriteHeadByte(StatCounterUplayerMessageTag)
	}

	lb.lower.OnUpperData(context)
}

// OnLowerData ...
func (lb *LoadBalancer) OnLowerData(context Context) {
	ub := context.GetBuffer()
	if ub == nil {
		return
	}

	tag, err := ub.ReadByte()
	if err != nil {
		return
	}

	if lb.enable {
		if tag != LoadBalancerUplayerMessageTag {
			fmt.Println("LoadBalancer: OnUpperData: todo")
			return
		}
	}

	lb.upper.OnLowerData(context)
}
