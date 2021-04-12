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

// OnUpperPush ...
func (frm *FrameDecoder) OnUpperPush(context Context) {
	if frm.enable {
		fmt.Println("FrameDecoder OnUpperPush")
	}

	frm.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (frm *FrameDecoder) OnLowerPush(context Context) {
	if frm.enable {
		fmt.Println("FrameDecoder OnLowerPush")
	}

	frm.upperDataProcessor.OnLowerPush(context)
}
