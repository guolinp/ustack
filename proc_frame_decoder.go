// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// FrameDecoder ...
type FrameDecoder struct {
	Base
}

// NewFrameDecoder ...
func NewFrameDecoder() DataProcessor {
	frm := &FrameDecoder{
		NewBaseInstance("FrameDecoder"),
	}
	return frm.Base.SetWhere(frm)
}

// OnUpperData ...
func (frm *FrameDecoder) OnUpperData(context Context) {
	if frm.enable {
		fmt.Println("FrameDecoder: OnUpperData: todo ")
	}

	frm.lower.OnUpperData(context)
}

// OnLowerData ...
func (frm *FrameDecoder) OnLowerData(context Context) {
	if frm.enable {
		fmt.Println("FrameDecoder: OnLowerData: todo")
	}

	frm.upper.OnLowerData(context)
}
