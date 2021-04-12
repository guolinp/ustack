// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// FlowController ...
type FlowController struct {
	Base
}

// NewFlowController ...
func NewFlowController() DataProcessor {
	fc := &FlowController{
		NewBaseInstance("FlowController"),
	}
	return fc.Base.SetWhere(fc)
}

// OnUpperPush ...
func (fc *FlowController) OnUpperPush(context Context) {
	if fc.enable {
		fmt.Println("FlowController: OnUpperPush: todo")
	}

	fc.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (fc *FlowController) OnLowerPush(context Context) {
	if fc.enable {
		fmt.Println("FlowController: OnLowerPush: todo")
	}

	fc.upperDataProcessor.OnLowerPush(context)
}
