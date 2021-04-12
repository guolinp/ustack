// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// StatCounter ...
type StatCounter struct {
	Base
	txCounter int
	rxCounter int
}

// NewStatCounter ...
func NewStatCounter() DataProcessor {
	sc := &StatCounter{
		Base:      NewBaseInstance("StatCounter"),
		txCounter: 0,
		rxCounter: 0,
	}
	return sc.Base.SetWhere(sc)
}

// OnUpperPush ...
func (sc *StatCounter) OnUpperPush(context Context) {
	if sc.enable {
		sc.txCounter++
		fmt.Println("StatCounter: txCounter:", sc.txCounter)
	}

	sc.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (sc *StatCounter) OnLowerPush(context Context) {
	if sc.enable {
		sc.rxCounter++
		fmt.Println("StatCounter: rxCounter:", sc.rxCounter)
	}

	sc.upperDataProcessor.OnLowerPush(context)
}
