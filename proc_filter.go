// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

type FilterFn func(context Context, toUpper bool) bool

// Filter ...
type Filter struct {
	ProcBase
	filterFn  []FilterFn
	txCounter uint64
	rxCounter uint64
}

// NewFilter ...
func NewFilter(filterFn ...FilterFn) DataProcessor {
	filter := &Filter{
		ProcBase:  NewProcBaseInstance("Filter"),
		filterFn:  make([]FilterFn, 1),
		txCounter: 0,
		rxCounter: 0,
	}

	filter.filterFn = append(filter.filterFn, filterFn...)

	return filter.ProcBase.SetWhere(filter)
}

// doFilter ...
func (filter *Filter) doFilter(context Context, toUpper bool) bool {
	for _, fn := range filter.filterFn {
		if fn(context, toUpper) {
			return true
		}
	}
	return false
}

// OnUpperData ...
func (filter *Filter) OnUpperData(context Context) {
	if filter.enable {
		if !filter.doFilter(context, false) {
			filter.txCounter++
			return
		}
	}

	filter.lower.OnUpperData(context)
}

// OnLowerData ...
func (filter *Filter) OnLowerData(context Context) {
	if filter.enable {
		if !filter.doFilter(context, true) {
			filter.rxCounter++
			return
		}
	}

	filter.upper.OnLowerData(context)
}
