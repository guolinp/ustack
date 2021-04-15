// Copyright 2021 The godevsig Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// UBuf is designed to manage the data between the various protocol
// layers of the communication protocol stack. When the message is
// constructed or parsed, data can be read and written from the head
// or tail, thereby avoiding frequent creation and deletion of caches
// or data movement. It also supports both big and little endian
// reading and writing.

// UBuf Format:
//
//      /------ reserved ------\
//     |                        |
//     +------------------------+------------------------+-----------------------+
//     | head space (to write)  | data space (to read)   | free space (to write) |
//     +------------------------+------------------------+-----------------------+
//     |                        ^                        ^                       ^
//     |                        |                        |                       |
//     |     WriteHeadXxx <---- | ----> ReadXxx          | ----> WriteXxx        |
//     |                        |                        |                       |
//     |                   readerIndex              writerIndex                  |
//     |                                                                         |
//      \----------------------------- capacity --------------------------------/
//

package ustack

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
)

// UBuf is struct to manage buffer
type UBuf struct {
	reserved    int
	readerIndex int
	writerIndex int
	data        []byte
}

// AllocWithHeadReserved returns UBuf instance or panic if invalid input given
//   capacity: the max bytes that Ubuf manages
//   reserved: bytes number of reserved in head space
func UBufAllocWithHeadReserved(capacity int, reserved int) *UBuf {
	if capacity <= 0 || reserved > capacity {
		log.Panicf("UBuf bad intput: capacity: %d, reserved: %d\n", capacity, reserved)
	}

	return &UBuf{
		reserved:    reserved,
		readerIndex: reserved,
		writerIndex: reserved,
		data:        make([]byte, capacity),
	}
}

// Alloc is shortcut version of AllocWithHeadReserved
// Do not reserve room in head space
func UBufAlloc(capacity int) *UBuf {
	return UBufAllocWithHeadReserved(capacity, 0)
}

// Reset reinitializes the UBuf, data will be lost
func (ub *UBuf) Reset() {
	ub.readerIndex = ub.reserved
	ub.writerIndex = ub.reserved
}

// Capacity returns the capacity
func (ub *UBuf) Capacity() int {
	return cap(ub.data)
}

// ReadableLength returns the length of readable data
func (ub *UBuf) ReadableLength() int {
	return ub.writerIndex - ub.readerIndex
}

// HeadWritableLength returns the length of writable data in head space
func (ub *UBuf) HeadWritableLength() int {
	return ub.readerIndex
}

// TailWritableLength returns the length of writable data in free space
func (ub *UBuf) TailWritableLength() int {
	return cap(ub.data) - ub.writerIndex
}

// Peek fills a byte slice with readable data, returns length of filled data or error
func (ub *UBuf) Peek(p []byte) (int, error) {
	toPeek := ub.ReadableLength()
	if toPeek <= 0 {
		return 0, errors.New("UBuf is not enough value to peek")
	}

	if toPeek > len(p) {
		toPeek = len(p)
	}

	toPeek = copy(p[:toPeek], ub.data[ub.readerIndex:])

	if toPeek < len(p) {
		return toPeek, errors.New("UBuf has not more data to peek")
	}

	return toPeek, nil
}

// PeekByte returns one byte with readable data or error
func (ub *UBuf) PeekByte() (byte, error) {
	if ub.ReadableLength() < 1 {
		return 0, errors.New("UBuf is not enough value to peek")
	}
	return ub.data[ub.readerIndex], nil
}

// PeekU16 returns a uint16 value with readable data or error
func (ub *UBuf) PeekU16() (uint16, error) {
	return ub.PeekU16BE()
}

// PeekU32 returns a uint32 value with readable data or error
func (ub *UBuf) PeekU32() (uint32, error) {
	return ub.PeekU32BE()
}

// PeekU64 returns a uint64 value with readable data or error
func (ub *UBuf) PeekU64() (uint64, error) {
	return ub.PeekU64BE()
}

// PeekU16BE returns a big endian uint16 value with readable data or error
func (ub *UBuf) PeekU16BE() (uint16, error) {
	if ub.ReadableLength() < 2 {
		return 0, errors.New("UBuf is not enough value to peek")
	}
	return binary.BigEndian.Uint16(ub.data[ub.readerIndex : ub.readerIndex+2]), nil
}

// PeekU32BE returns a big endian uint32 value with readable data or error
func (ub *UBuf) PeekU32BE() (uint32, error) {
	if ub.ReadableLength() < 4 {
		return 0, errors.New("UBuf is not enough value to peek")
	}
	return binary.BigEndian.Uint32(ub.data[ub.readerIndex : ub.readerIndex+4]), nil
}

// PeekU64BE returns a big endian uint64 value with readable data or error
func (ub *UBuf) PeekU64BE() (uint64, error) {
	if ub.ReadableLength() < 8 {
		return 0, errors.New("UBuf is not enough value to peek")
	}
	return binary.BigEndian.Uint64(ub.data[ub.readerIndex : ub.readerIndex+8]), nil
}

// PeekU16LE returns a little endian uint16 value with readable data or error
func (ub *UBuf) PeekU16LE() (uint16, error) {
	if ub.ReadableLength() < 2 {
		return 0, errors.New("UBuf is not enough value to peek")
	}
	return binary.LittleEndian.Uint16(ub.data[ub.readerIndex : ub.readerIndex+2]), nil
}

// PeekU32LE returns a little endian uint32 value with readable data or error
func (ub *UBuf) PeekU32LE() (uint32, error) {
	if ub.ReadableLength() < 4 {
		return 0, errors.New("UBuf is not enough value to peek")
	}
	return binary.LittleEndian.Uint32(ub.data[ub.readerIndex : ub.readerIndex+4]), nil
}

// PeekU64LE returns a little endian uint64 value with readable data or error
func (ub *UBuf) PeekU64LE() (uint64, error) {
	if ub.ReadableLength() < 8 {
		return 0, errors.New("UBuf is not enough value to peek")
	}
	return binary.LittleEndian.Uint64(ub.data[ub.readerIndex : ub.readerIndex+8]), nil
}

// WriteByte implements io.ByteWriter interface
func (ub *UBuf) WriteByte(c byte) error {
	if ub.TailWritableLength() <= 0 {
		return errors.New("UBuf is full")
	}

	ub.data[ub.writerIndex] = c
	ub.writerIndex++

	return nil
}

// Write implements io.Writer interface
func (ub *UBuf) Write(p []byte) (n int, err error) {
	toWrite := ub.TailWritableLength()
	if toWrite <= 0 {
		return 0, errors.New("UBuf is full")
	}

	if toWrite > len(p) {
		toWrite = len(p)
	}

	toWrite = copy(ub.data[ub.writerIndex:], p[:toWrite])
	ub.writerIndex += toWrite

	if toWrite < len(p) {
		return toWrite, errors.New("UBuf is wrote partial data")
	}

	return toWrite, nil
}

// WriteTo implements io.WriterTo interface
func (ub *UBuf) WriteTo(w io.Writer) (n int64, err error) {
	if ub.ReadableLength() <= 0 {
		return 0, nil
	}

	written, err := w.Write(ub.data[ub.readerIndex:ub.writerIndex])
	if err != nil {
		return int64(written), err
	}

	ub.readerIndex += written

	return int64(written), nil
}

// WriteU16 writes uint16 data into data space
func (ub *UBuf) WriteU16(value uint16) error {
	return ub.WriteU16BE(value)
}

// WriteU32 writes uint32 data into data space
func (ub *UBuf) WriteU32(value uint32) error {
	return ub.WriteU32BE(value)
}

// WriteU64 writes uint64 data into data space
func (ub *UBuf) WriteU64(value uint64) error {
	return ub.WriteU64BE(value)
}

// WriteU16BE writes big endian uint16 data into data space
func (ub *UBuf) WriteU16BE(value uint16) error {
	return binary.Write(ub, binary.BigEndian, value)
}

// WriteU32BE writes big endian uint32 data into data space
func (ub *UBuf) WriteU32BE(value uint32) error {
	return binary.Write(ub, binary.BigEndian, value)
}

// WriteU64BE writes big endian uint64 data into data space
func (ub *UBuf) WriteU64BE(value uint64) error {
	return binary.Write(ub, binary.BigEndian, value)
}

// WriteU16LE writes little endian uint16 data into data space
func (ub *UBuf) WriteU16LE(value uint16) error {
	return binary.Write(ub, binary.LittleEndian, value)
}

// WriteU32LE writes little endian uint32 data into data space
func (ub *UBuf) WriteU32LE(value uint32) error {
	return binary.Write(ub, binary.LittleEndian, value)
}

// WriteU64LE writes little endian uint64 data into data space
func (ub *UBuf) WriteU64LE(value uint64) error {
	return binary.Write(ub, binary.LittleEndian, value)
}

// ReadByte implements io.ByteReader interface
func (ub *UBuf) ReadByte() (byte, error) {
	if ub.ReadableLength() <= 0 {
		return 0, errors.New("UBuf is empty")
	}

	b := ub.data[ub.readerIndex]
	ub.readerIndex++

	return b, nil
}

// Read implements io.Reader interface
func (ub *UBuf) Read(p []byte) (n int, err error) {
	toRead := ub.ReadableLength()
	if toRead <= 0 {
		return 0, nil
	}

	if toRead > len(p) {
		toRead = len(p)
	}

	toRead = copy(p[:toRead], ub.data[ub.readerIndex:])
	ub.readerIndex += toRead

	if toRead < len(p) {
		return toRead, errors.New("UBuf has not more data to read")
	}

	return toRead, nil
}

// ReadFrom implements io.ReaderFrom interface
func (ub *UBuf) ReadFrom(r io.Reader) (n int64, err error) {
	if ub.TailWritableLength() <= 0 {
		return 0, nil
	}

	length, err := r.Read(ub.data[ub.writerIndex:])
	if err != nil {
		return 0, err
	}

	ub.writerIndex += length

	return int64(length), nil
}

// ReadU16 returns uint16 data or error
func (ub *UBuf) ReadU16() (uint16, error) {
	return ub.ReadU16BE()
}

// ReadU32 returns uint32 data or error
func (ub *UBuf) ReadU32() (uint32, error) {
	return ub.ReadU32BE()
}

// ReadU64 returns uint64 data or error
func (ub *UBuf) ReadU64() (uint64, error) {
	return ub.ReadU64BE()
}

// ReadU16BE returns big endian uint16 data or error
func (ub *UBuf) ReadU16BE() (uint16, error) {
	var value uint16
	err := binary.Read(ub, binary.BigEndian, &value)
	return value, err
}

// ReadU32BE returns big endian uint32 data or error
func (ub *UBuf) ReadU32BE() (uint32, error) {
	var value uint32
	err := binary.Read(ub, binary.BigEndian, &value)
	return value, err
}

// ReadU64BE returns big endian uint64 data or error
func (ub *UBuf) ReadU64BE() (uint64, error) {
	var value uint64
	err := binary.Read(ub, binary.BigEndian, &value)
	return value, err
}

// ReadU16LE returns little endian uint16 data or error
func (ub *UBuf) ReadU16LE() (uint16, error) {
	var value uint16
	err := binary.Read(ub, binary.LittleEndian, &value)
	return value, err
}

// ReadU32LE returns little endian uint32 data or error
func (ub *UBuf) ReadU32LE() (uint32, error) {
	var value uint32
	err := binary.Read(ub, binary.LittleEndian, &value)
	return value, err
}

// ReadU64LE returns little endian uint64 data or error
func (ub *UBuf) ReadU64LE() (uint64, error) {
	var value uint64
	err := binary.Read(ub, binary.LittleEndian, &value)
	return value, err
}

// writeHead is a helper function, it writes data into head space
func (ub *UBuf) writeHead(dataSize int, fn func() error) error {
	if ub.readerIndex < dataSize {
		return errors.New("UBuf Head room is not enouth")
	}

	// all the writes are from writerIndex
	// backup the writerIndex first
	oldWriterIndex := ub.writerIndex

	// move writerIndex to new position where we shall write data
	ub.writerIndex = ub.readerIndex - dataSize
	err := fn()
	if err == nil {
		// write success, update readerIndex
		ub.readerIndex -= dataSize
	}

	// restore the writerIndex
	ub.writerIndex = oldWriterIndex

	return err
}

// WriteHeadByte writes one byte into head space
func (ub *UBuf) WriteHeadByte(value byte) error {
	return ub.writeHead(1, func() error { return ub.WriteByte(value) })
}

// WriteHeadU16 writes uint16 data into head space
func (ub *UBuf) WriteHeadU16(value uint16) error {
	return ub.WriteHeadU16BE(value)
}

// WriteHeadU32 writes uint32 data into head space
func (ub *UBuf) WriteHeadU32(value uint32) error {
	return ub.WriteHeadU32BE(value)
}

// WriteHeadU64 writes uint64 data into head space
func (ub *UBuf) WriteHeadU64(value uint64) error {
	return ub.WriteHeadU64BE(value)
}

// WriteHeadU16BE writes big endian uint16 data into head space
func (ub *UBuf) WriteHeadU16BE(value uint16) error {
	return ub.writeHead(2, func() error { return ub.WriteU16BE(value) })
}

// WriteHeadU32BE writes big endian uint32 data into head space
func (ub *UBuf) WriteHeadU32BE(value uint32) error {
	return ub.writeHead(4, func() error { return ub.WriteU32BE(value) })
}

// WriteHeadU64BE writes big endian uint64 data into head space
func (ub *UBuf) WriteHeadU64BE(value uint64) error {
	return ub.writeHead(8, func() error { return ub.WriteU64BE(value) })
}

// WriteHeadU16LE writes little endian uint16 data into head space
func (ub *UBuf) WriteHeadU16LE(value uint16) error {
	return ub.writeHead(2, func() error { return ub.WriteU16LE(value) })
}

// WriteHeadU32LE writes little endian uint32 data into head space
func (ub *UBuf) WriteHeadU32LE(value uint32) error {
	return ub.writeHead(4, func() error { return ub.WriteU32LE(value) })
}

// WriteHeadU64LE writes little endian uint64 data into head space
func (ub *UBuf) WriteHeadU64LE(value uint64) error {
	return ub.writeHead(8, func() error { return ub.WriteU64LE(value) })
}
