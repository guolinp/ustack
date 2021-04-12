// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Forwarder ...
type Forwarder struct {
	Base
}

// NewForwarder ...
func NewForwarder() DataProcessor {
	fwd := &Forwarder{
		NewBaseInstance("Forwarder"),
	}
	return fwd.Base.SetWhere(fwd)
}

// OnUpperPush ...
func (fwd *Forwarder) OnUpperPush(context Context) {
	if fwd.enable {
		fmt.Println("Forwarder OnUpperPush")
	}

	fwd.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (fwd *Forwarder) OnLowerPush(context Context) {
	if fwd.enable {
		fmt.Println("Forwarder OnLowerPush")
	}

	fwd.upperDataProcessor.OnLowerPush(context)
}
