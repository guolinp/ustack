// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Echo ...
type Echo struct {
	ProcBase
}

// NewEcho ...
func NewEcho() DataProcessor {
	echo := &Echo{
		NewProcBaseInstance("Echo"),
	}
	return echo.ProcBase.SetWhere(echo)
}

// OnUpperData ...
func (echo *Echo) OnUpperData(context Context) {
	echo.lower.OnUpperData(context)
}

// OnLowerData ...
func (echo *Echo) OnLowerData(context Context) {
	if echo.enable {
		fmt.Println("Echo: send back the lowlayer data")
		echo.lower.OnUpperData(context)
	} else {
		echo.upper.OnLowerData(context)
	}
}
