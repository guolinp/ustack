// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Discarder ...
type Discarder struct {
	ProcBase
}

// NewDiscarder ...
func NewDiscarder() DataProcessor {
	dis := &Discarder{
		NewProcBaseInstance("Discarder"),
	}
	return dis.ProcBase.SetWhere(dis)
}

// OnUpperData ...
func (dis *Discarder) OnUpperData(context Context) {
	dis.lower.OnUpperData(context)
}

// OnLowerData ...
func (dis *Discarder) OnLowerData(context Context) {
	if dis.enable {
		fmt.Println("Discarder: drop the lowlayer data")
	} else {
		dis.upper.OnLowerData(context)
	}
}
