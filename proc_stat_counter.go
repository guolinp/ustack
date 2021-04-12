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

// OnUpperData ...
func (sc *StatCounter) OnUpperData(context Context) {
	if sc.enable {
		sc.txCounter++
		fmt.Println("StatCounter: txCounter:", sc.txCounter)
	}

	sc.lower.OnUpperData(context)
}

// OnLowerData ...
func (sc *StatCounter) OnLowerData(context Context) {
	if sc.enable {
		sc.rxCounter++
		fmt.Println("StatCounter: rxCounter:", sc.rxCounter)
	}

	sc.upper.OnLowerData(context)
}
