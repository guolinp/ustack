// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// StatCounter ...
type StatCounter struct {
	Base
}

// NewStatCounter ...
func NewStatCounter() DataProcessor {
	sc := &StatCounter{
		NewBaseInstance("StatCounter"),
	}
	return sc.Base.SetWhere(sc)
}

// OnUpperPush ...
func (sc *StatCounter) OnUpperPush(context Context) {
	if sc.enable {
		fmt.Println("StatCounter OnUpperPush")
	}

	sc.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (sc *StatCounter) OnLowerPush(context Context) {
	if sc.enable {
		fmt.Println("StatCounter OnLowerPush")
	}

	sc.upperDataProcessor.OnLowerPush(context)
}
