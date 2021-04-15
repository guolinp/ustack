// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Forwarder ...
type Forwarder struct {
	ProcBase
}

// NewForwarder ...
func NewForwarder() DataProcessor {
	fwd := &Forwarder{
		NewProcBaseInstance("Forwarder"),
	}
	return fwd.ProcBase.SetWhere(fwd)
}

// OnUpperData ...
func (fwd *Forwarder) OnUpperData(context Context) {
	if fwd.enable {
		fmt.Println("Forwarder: forward uplayer data, todo")
	}

	fwd.lower.OnUpperData(context)
}

// OnLowerData ...
func (fwd *Forwarder) OnLowerData(context Context) {
	if fwd.enable {
		fmt.Println("Forwarder: forward lowlayer data, todo")
	}

	fwd.upper.OnLowerData(context)
}
