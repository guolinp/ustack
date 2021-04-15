// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Filter ...
type Filter struct {
	ProcBase
}

// NewFilter ...
func NewFilter() DataProcessor {
	filter := &Filter{
		NewProcBaseInstance("Filter"),
	}
	return filter.ProcBase.SetWhere(filter)
}

// OnUpperData ...
func (filter *Filter) OnUpperData(context Context) {
	if filter.enable {
		fmt.Println("Filter: drop uplayer data")
	} else {
		filter.lower.OnUpperData(context)
	}
}

// OnLowerData ...
func (filter *Filter) OnLowerData(context Context) {
	if filter.enable {
		fmt.Println("Filter: drop lowlayer data")
	} else {
		filter.upper.OnLowerData(context)
	}
}
