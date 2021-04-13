// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

func OptionParseByte(option interface{}, defaultValue byte) (value byte, exits bool) {
	if option != nil {
		value, ok := option.(byte)
		if ok {
			return value, true
		}
	}
	return defaultValue, false
}

func OptionParseInt(option interface{}, defaultValue int) (value int, exits bool) {
	if option != nil {
		value, ok := option.(int)
		if ok {
			return value, true
		}
	}
	return defaultValue, false
}

func OptionParseBool(option interface{}, defaultValue bool) (value bool, exits bool) {
	if option != nil {
		value, ok := option.(bool)
		if ok {
			return value, true
		}
	}
	return defaultValue, false
}

func OptionParseByteSlice(option interface{}) (value []byte, ok bool) {
	if option != nil {
		value, ok := option.([]byte)
		if ok {
			return value, true
		}
	}
	return nil, false
}
