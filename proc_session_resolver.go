// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// SessionResolver ...
type SessionResolver struct {
	Base
}

// NewSessionResolver ...
func NewSessionResolver() DataProcessor {
	sr := &SessionResolver{
		NewBaseInstance("Session"),
	}
	return sr.Base.SetWhere(sr)
}

// OnUpperPush ...
func (sr *SessionResolver) OnUpperPush(context Context) {
	if sr.enable {
		s, ok := context.GetOption("session")

		var session byte = 0
		if ok {
			session, ok = s.(byte)
			if !ok {
				return
			}
		}

		ub := context.GetBuffer()
		if ub == nil {
			return
		}

		err := ub.WriteHeadByte(session)
		if err != nil {
			return
		}
	}

	sr.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (sr *SessionResolver) OnLowerPush(context Context) {
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

	sr.upperDataProcessor.OnLowerPush(context)
}
