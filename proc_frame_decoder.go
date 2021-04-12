// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ustack

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
		ub.WriteHeadU16BE(uint16(ub.ReadableLength()))
	}

	frm.lower.OnUpperData(context)
}

// OnLowerData ...
func (frm *FrameDecoder) OnLowerData(context Context) {
	if frm.enable {
		ub := context.GetBuffer()
		if ub == nil {
			return
		}

		if frm.cache == nil {
			frm.cache = ub
		} else {
			frm.cache.ReadFrom(ub)
		}

		length := frm.cache.ReadableLength()
		if length < FrameLengthFieldSizeInByte {
			// wait for more data
			return
		}

		frameLength, err := frm.cache.PeekU16BE()
		if err != nil {
			// bad buffer, discard it
			frm.cache = nil
			return
		}

		// there must have a complete frame
		if length >= int(frameLength) {
			newUbuf := UBufAlloc(int(frameLength))

			// drop size-field-data by dummy reading
			frm.cache.ReadU16BE()

			// fill the new buffer for uplayer
			newUbuf.ReadFrom(frm.cache)

			// cache is empty, free the reference
			if frm.cache.ReadableLength() <= 0 {
				frm.cache = nil
			}

			context.SetBuffer(newUbuf)

			frm.upper.OnLowerData(context)
		}
	} else {
		frm.upper.OnLowerData(context)
	}
}
