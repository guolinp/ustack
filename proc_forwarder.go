// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

type ForwardFn func(context Context, toUpper bool) bool

// Forwarder ...
type Forwarder struct {
	ProcBase
	forwardFn []ForwardFn
	txCounter uint64
	rxCounter uint64
}

// NewForwarder ...
func NewForwarder(forwardFn ...ForwardFn) DataProcessor {
	fwd := &Forwarder{
		ProcBase:  NewProcBaseInstance("Forwarder"),
		forwardFn: make([]ForwardFn, 1),
		txCounter: 0,
		rxCounter: 0,
	}

	fwd.forwardFn = append(fwd.forwardFn, forwardFn...)

	return fwd.ProcBase.SetWhere(fwd)
}

// doForward ...
func (fwd *Forwarder) doForward(context Context, toUpper bool) bool {
	for _, fn := range fwd.forwardFn {
		if fn(context, toUpper) {
			return true
		}
	}
	return false
}

// OnUpperData ...
func (fwd *Forwarder) OnUpperData(context Context) {
	if fwd.enable {
		if fwd.doForward(context, false) {
			fwd.txCounter++
			return
		}
	}

	fwd.lower.OnUpperData(context)
}

// OnLowerData ...
func (fwd *Forwarder) OnLowerData(context Context) {
	if fwd.enable {
		if fwd.doForward(context, true) {
			fwd.rxCounter++
			return
		}
	}

	fwd.upper.OnLowerData(context)
}
