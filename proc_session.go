// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// SessionResolver ...
type SessionResolver struct {
	name               string
	enable             bool
	options            map[string]interface{}
	upperDataProcessor DataProcessor
	lowerDataProcessor DataProcessor
}

// NewSessionResolver ...
func NewSessionResolver() DataProcessor {
	return &SessionResolver{
		name:    "SessionResolver",
		enable:  true,
		options: make(map[string]interface{}),
	}
}

// SetName ...
func (pr *SessionResolver) SetName(name string) DataProcessor {
	pr.name = name
	return pr
}

// GetName ...
func (pr *SessionResolver) GetName() string {
	return pr.name
}

// SetOption ...
func (pr *SessionResolver) SetOption(name string, value interface{}) DataProcessor {
	pr.options[name] = value
	return pr
}

// GetOption ...
func (pr *SessionResolver) GetOption(name string) interface{} {
	if value, ok := pr.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable ...
func (pr *SessionResolver) SetEnable(enable bool) DataProcessor {
	pr.enable = enable
	return pr
}

// ForServer ...
func (pr *SessionResolver) ForServer(forServer bool) DataProcessor {
	return pr
}

// SetUStack ...
func (pr *SessionResolver) SetUStack(u UStack) DataProcessor {
	return pr
}

// SetUpperDataProcessor ...
func (pr *SessionResolver) SetUpperDataProcessor(dp DataProcessor) DataProcessor {
	pr.upperDataProcessor = dp
	return pr
}

// SetLowerDataProcessor ...
func (pr *SessionResolver) SetLowerDataProcessor(dp DataProcessor) DataProcessor {
	pr.lowerDataProcessor = dp
	return pr
}

// OnUpperPush ...
func (pr *SessionResolver) OnUpperPush(context Context) {
	if pr.enable {
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

	pr.lowerDataProcessor.OnUpperPush(context)
}

// OnLowerPush ...
func (pr *SessionResolver) OnLowerPush(context Context) {
	if pr.enable {
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

	pr.upperDataProcessor.OnLowerPush(context)
}

// OnEvent ...
func (pr *SessionResolver) OnEvent(event Event) {
}

// Run ...
func (pr *SessionResolver) Run() DataProcessor {
	return pr
}
