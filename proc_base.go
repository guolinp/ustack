// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

// Base as a special data processor is used to manage and maintain
// the common operations of all data processor. it is usually embedded
// in other data processor and should NOT be used directly
type Base struct {
	where     DataProcessor
	name      string
	enable    bool
	ustack    UStack
	forServer bool
	options   map[string]interface{}
	upper     DataProcessor
	lower     DataProcessor
}

// NewBaseInstance returns a new instance
func NewBaseInstance(name string) Base {
	base := Base{
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

// NewBase return a new instance that meets for DataProcessor interface
func NewBase() DataProcessor {
	base := NewBaseInstance("Base")
	return &base
}

// SetWhere where the base is embedded in
func (base *Base) SetWhere(where DataProcessor) DataProcessor {
	base.where = where
	return base.where
}

// SetName set the name
func (base *Base) SetName(name string) DataProcessor {
	base.name = name
	return base.where
}

// GetName returns name
func (base *Base) GetName() string {
	return base.name
}

// SetOption set the options
//     name: option name
//     value: option value
func (base *Base) SetOption(name string, value interface{}) DataProcessor {
	base.options[name] = value
	return base.where
}

// GetOption returns the option vaule of given name
// return nil if the option does not exist
func (base *Base) GetOption(name string) interface{} {
	if value, ok := base.options[name]; ok {
		return value
	}
	return nil
}

// SetEnable enable(true) or disable(false) the DataProcessor
func (base *Base) SetEnable(enable bool) DataProcessor {
	base.enable = enable
	return base.where
}

// ForServer set
func (base *Base) ForServer(forServer bool) DataProcessor {
	base.forServer = forServer
	return base.where
}

// SetUStack set the UStack instance
func (base *Base) SetUStack(ustack UStack) DataProcessor {
	base.ustack = ustack
	return base.where
}

// SetUpper set upper data processor instance
func (base *Base) SetUpper(upper DataProcessor) DataProcessor {
	base.upper = upper
	return base.where
}

// SetLower set lower data processor instance
func (base *Base) SetLower(lower DataProcessor) DataProcessor {
	base.lower = lower
	return base.where
}

// OnUpperData is called when upper layer sending data
func (base *Base) OnUpperData(context Context) {
}

// OnLowerData is called when lower layer received data
func (base *Base) OnLowerData(context Context) {
}

// OnEvent is called when any event hanppen
func (base *Base) OnEvent(event Event) {
}

// Run starts the data processor
func (base *Base) Run() DataProcessor {
	return base.where
}
