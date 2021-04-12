// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "fmt"

// Echo ...
type Echo struct {
	Base
}

// NewEcho ...
func NewEcho() DataProcessor {
	echo := &Echo{
		NewBaseInstance("Echo"),
	}
	return echo.Base.SetWhere(echo)
}

// OnUpperPush ...
func (echo *Echo) OnUpperPush(context Context) {
	echo.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (echo *Echo) OnLowerPush(context Context) {
	if echo.enable {
		fmt.Println("Echo: send back the lowlayer data")
		echo.lowerDataProcessor.OnUpperPush(context)
	} else {
		echo.upperDataProcessor.OnLowerPush(context)
	}
}
