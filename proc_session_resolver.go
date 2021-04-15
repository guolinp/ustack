// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// SessionResolver ...
type SessionResolver struct {
	ProcBase
}

// NewSessionResolver ...
func NewSessionResolver() DataProcessor {
	sr := &SessionResolver{
		NewProcBaseInstance("SessionResolver"),
	}
	return sr.ProcBase.SetWhere(sr)
}

// GetOverhead ...
func (sr *SessionResolver) GetOverhead() int {
	// for saving session byte
	return 1
}

// OnUpperData ...
func (sr *SessionResolver) OnUpperData(context Context) {
	if sr.enable {
		ub := context.GetBuffer()
		if ub == nil {
			return
		}

		session, _ := OptionParseByte(context.GetOption("session"), 0)

		err := ub.WriteHeadByte(session)
		if err != nil {
			return
		}
	}

	sr.lower.OnUpperData(context)
}

// OnLowerData ...
func (sr *SessionResolver) OnLowerData(context Context) {
	if sr.enable {
		ub := context.GetBuffer()
		if ub == nil {
			return
		}

		session, err := ub.ReadByte()
		if err != nil {
			return
		}

		context.SetOption("session", session)
	}

	sr.upper.OnLowerData(context)
}
