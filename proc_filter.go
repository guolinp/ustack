// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Filter ...
type Filter struct {
	Base
}

// NewFilter ...
func NewFilter() DataProcessor {
	filter := &Filter{
		NewBaseInstance("Filter"),
	}
	return filter.Base.SetWhere(filter)
}

// OnUpperPush ...
func (filter *Filter) OnUpperPush(context Context) {
	if filter.enable {
		fmt.Println("Filter: drop uplayer data")
	} else {
		filter.lowerDataProcessor.OnUpperPush(context)
	}
}

// OnLowerPush ...
func (filter *Filter) OnLowerPush(context Context) {
	if filter.enable {
		fmt.Println("Filter: drop lowlayer data")
	} else {
		filter.upperDataProcessor.OnLowerPush(context)
	}
}
