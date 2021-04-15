// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// FlowController ...
type FlowController struct {
	ProcBase
}

// NewFlowController ...
func NewFlowController() DataProcessor {
	fc := &FlowController{
		NewProcBaseInstance("FlowController"),
	}
	return fc.ProcBase.SetWhere(fc)
}

// OnUpperData ...
func (fc *FlowController) OnUpperData(context Context) {
	if fc.enable {
		fmt.Println("FlowController: OnUpperData: todo")
	}

	fc.lower.OnUpperData(context)
}

// OnLowerData ...
func (fc *FlowController) OnLowerData(context Context) {
	if fc.enable {
		fmt.Println("FlowController: OnLowerData: todo")
	}

	fc.upper.OnLowerData(context)
}
