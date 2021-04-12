// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Discarder ...
type Discarder struct {
	Base
}

// NewDiscarder ...
func NewDiscarder() DataProcessor {
	dis := &Discarder{
		NewBaseInstance("Discarder"),
	}
	return dis.Base.SetWhere(dis)
}

// OnUpperPush ...
func (dis *Discarder) OnUpperPush(context Context) {
	if dis.enable {
		fmt.Println("Discarder OnUpperPush")
	}

	dis.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (dis *Discarder) OnLowerPush(context Context) {
	if dis.enable {
		fmt.Println("Discarder OnLowerPush")
	}

	dis.upperDataProcessor.OnLowerPush(context)
}
