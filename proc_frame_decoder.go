// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

import "io"

const FrameLengthFieldSizeInByte int = 2

// FrameDecoder ...
type FrameDecoder struct {
	Base
	cache *UBuf
}

// NewFrameDecoder ...
func NewFrameDecoder() DataProcessor {
	frm := &FrameDecoder{
		Base:  NewBaseInstance("FrameDecoder"),
		cache: nil,
	}
	return frm.Base.SetWhere(frm)
}

// GetOverhead returns the overhead
func (frm *FrameDecoder) GetOverhead() int {
	return FrameLengthFieldSizeInByte
}

// OnUpperData ...
func (frm *FrameDecoder) OnUpperData(context Context) {
	if frm.enable {
		ub := context.GetBuffer()
		if ub == nil {
			return
		}

		// add frame length in head space
		ub.WriteHeadU16BE(uint16(ub.ReadableLength()))
	}

	frm.lower.OnUpperData(context)
}

// handleCurrentData ...
func (frm *FrameDecoder) handleCurrentData(context Context, ub *UBuf) {
	// if there is data in cache, put current data into cache
	if frm.cache.ReadableLength() > 0 {
		frm.cache.ReadFrom(ub)
		return
	}

	// cache is empty
	// handle as much as possiable with loop
	for {
		// very less data, cache the data
		if ub.ReadableLength() < FrameLengthFieldSizeInByte {
			frm.cache.ReadFrom(ub)
			return
		}

		expectedLength, err := ub.PeekU16BE()
		if err != nil {
			// bad buffer, discard it
			frm.cache.Reset()
			return
		}

		actuallyLength := uint16(ub.ReadableLength())

		// not a complete frame, cache the data
		if expectedLength > actuallyLength {
			frm.cache.ReadFrom(ub)
			return
		}

		// just one frame, need not to alloc new UBuf
		if expectedLength == actuallyLength {
			// drop size-field-data by dummy reading
			ub.ReadU16BE()

			context.SetBuffer(ub)
			frm.upper.OnLowerData(context)
			return
		}

		// here: expectedLength < actuallyLength
		// there must have at least one complete frame
		newUbuf := UBufAlloc(int(expectedLength))

		// drop size-field-data by dummy reading
		ub.ReadU16BE()

		// fill the new buffer for uplayer
		_, err = io.CopyN(newUbuf, ub, int64(expectedLength))
		if err != nil {
			// bad buffer, discard it
			frm.cache.Reset()
			return
		}

		context.SetBuffer(newUbuf)

		// invoke uplayer
		frm.upper.OnLowerData(context)
	}
}

// handleCachedData ...
func (frm *FrameDecoder) handleCachedData(context Context, ub *UBuf) {
	// handle as much as possiable with loop
	for {
		cachedLength := frm.cache.ReadableLength()

		if cachedLength == 0 {
			frm.cache.Reset()
			return
		}

		if cachedLength < FrameLengthFieldSizeInByte {
			// wait for more data
			return
		}

		expectedLength, err := frm.cache.PeekU16BE()
		if err != nil {
			// bad buffer, discard it
			frm.cache.Reset()
			return
		}

		if cachedLength < int(expectedLength) {
			// wait for more data
			return
		}

		// there must have at least one complete frame
		newUbuf := UBufAlloc(int(expectedLength))

		// drop size-field-data by dummy reading
		frm.cache.ReadU16BE()

		// fill the new buffer for uplayer
		_, err = io.CopyN(newUbuf, frm.cache, int64(expectedLength))
		if err != nil {
			// bad buffer, discard it
			frm.cache.Reset()
			return
		}

		context.SetBuffer(newUbuf)

		// invoke uplayer
		frm.upper.OnLowerData(context)
	}
}

// OnLowerData ...
func (frm *FrameDecoder) OnLowerData(context Context) {
	ub := context.GetBuffer()
	if ub == nil {
		return
	}

	if frm.enable {
		// handle current received data
		frm.handleCurrentData(context, ub)

		// handle history cached data
		frm.handleCachedData(context, ub)
	} else {
		frm.upper.OnLowerData(context)
	}
}

// Run ...
func (frm *FrameDecoder) Run() DataProcessor {
	// 2 * MTU is for the worst case
	frm.cache = UBufAlloc(2 * frm.ustack.GetMTU())
	return frm
}
