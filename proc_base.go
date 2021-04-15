// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// ProcBase as a special data processor is used to manage and maintain
// the common operations of all data processor. it is usually embedded
// in other data processor and should NOT be used directly
type ProcBase struct {
	where     DataProcessor
	name      string
	enable    bool
	ustack    UStack
	forServer bool
	options   map[string]interface{}
	upper     DataProcessor
	lower     DataProcessor
}

// NewProcBaseInstance returns a new instance
func NewProcBaseInstance(name string) ProcBase {
	base := ProcBase{
		name:      name,
		enable:    true,
		ustack:    nil,
		forServer: true,
		options:   make(map[string]interface{}),
	}
	// by default is itself
	base.where = &base
	return base
}

// NewProcBase return a new instance that meets for DataProcessor interface
func NewProcBase() DataProcessor {
	base := NewProcBaseInstance("ProcBase")
	return &base
}

// SetWhere where the base is embedded in
func (base *ProcBase) SetWhere(where DataProcessor) DataProcessor {
	base.where = where
	return base.where
}

// SetName set the name
func (base *ProcBase) SetName(name string) DataProcessor {
	base.name = name
	return base.where
}

// GetName returns name
func (base *ProcBase) GetName() string {
	return base.name
}

// GetOverhead returns the overhead
func (base *ProcBase) GetOverhead() int {
	return 0
}

// SetOption set the options
//     name: option name
//     value: option value
func (base *ProcBase) SetOption(name string, value interface{}) DataProcessor {
	base.options[name] = value
	return base.where
}

// GetOption returns the option vaule of given name
// return nil if the option does not exist
func (base *ProcBase) GetOption(name string) interface{} {
	if value, ok := base.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable enable(true) or disable(false) the DataProcessor
func (base *ProcBase) SetEnable(enable bool) DataProcessor {
	base.enable = enable
	return base.where
}

// ForServer set
func (base *ProcBase) ForServer(forServer bool) DataProcessor {
	base.forServer = forServer
	return base.where
}

// SetUStack set the UStack instance
func (base *ProcBase) SetUStack(ustack UStack) DataProcessor {
	base.ustack = ustack
	return base.where
}

// SetUpper set upper data processor instance
func (base *ProcBase) SetUpper(upper DataProcessor) DataProcessor {
	base.upper = upper
	return base.where
}

// SetLower set lower data processor instance
func (base *ProcBase) SetLower(lower DataProcessor) DataProcessor {
	base.lower = lower
	return base.where
}

// OnUpperData is called when upper layer sending data
func (base *ProcBase) OnUpperData(context Context) {
}

// OnLowerData is called when lower layer received data
func (base *ProcBase) OnLowerData(context Context) {
}

// OnEvent is called when any event hanppen
func (base *ProcBase) OnEvent(event Event) {
}

// Run starts the data processor
func (base *ProcBase) Run() DataProcessor {
	return base.where
}
