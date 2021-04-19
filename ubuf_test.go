package ustack

import (
	"bytes"
	"fmt"
	"testing"
)

const UBufCapacity = 64
const UBufReserved = 8
const UBufHeadSize = UBufReserved
const UBufDataSize = UBufCapacity - UBufReserved

func TestUBufAlloc(t *testing.T) {
	ub := UBufAlloc(UBufCapacity)
	if ub.Capacity() != UBufCapacity {
		t.Fatal("Unexpected capacity")
	}

	if ub.ReadableLength() != 0 {
		t.Fatal("Unexpected readable length")
	}

	if ub.HeadWritableLength() != 0 {
		t.Fatal("Unexpected head writable length")
	}

	if ub.TailWritableLength() != UBufCapacity {
		t.Fatal("Unexpected tail writable length")
	}
}

func TestUBufAllocWithHeadReserved(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufReserved)
	if ub.Capacity() != UBufCapacity {
		t.Fatal("Unexpected capacity")
	}

	if ub.ReadableLength() != 0 {
		t.Fatal("Unexpected readable length")
	}

	if ub.HeadWritableLength() != UBufReserved {
		t.Fatal("Unexpected head writable length")
	}

	if ub.TailWritableLength() != UBufCapacity-UBufReserved {
		t.Fatal("Unexpected tail writable length")
	}
}

func TestBadUBufAllocInvalidCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("The code did not panic")
		}
	}()

	UBufAlloc(-1)
}

func TestBadUBufAllocInvalidReserved(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("The code did not panic")
		}
	}()

	UBufAllocWithHeadReserved(4, 5)
}

func TestReadWriteByte(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufReserved)

	for i := 0; i < UBufDataSize; i++ {
		err := ub.WriteByte(1)
		if err != nil {
			t.Fatal("Unexpected write byte result")
		}

		if ub.TailWritableLength() != UBufDataSize-i-1 {
			t.Fatal("Unexpected writable length")
		}
	}

	err := ub.WriteByte(1)
	if err == nil {
		t.Fatal("Unexpected write byte result")
	}

	for i := 0; i < UBufDataSize; i++ {
		b, err := ub.ReadByte()
		if err != nil {
			t.Fatal("Unexpected read byte result")
		}

		if b != 1 {
			t.Fatal("Unexpected read byte value")
		}

		if ub.ReadableLength() != UBufDataSize-i-1 {
			t.Fatal("Unexpected readable length")
		}
	}

	if ub.TailWritableLength() != 0 {
		t.Fatal("Unexpected tail writable length")
	}

	if ub.ReadableLength() != 0 {
		t.Fatal("Unexpected readable length")
	}

	_, err = ub.ReadByte()
	if err == nil {
		t.Fatal("Unexpected read byte result")
	}
}

func TestReadWrite(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufReserved)

	writeData := make([]byte, UBufDataSize+1)
	readData := make([]byte, UBufDataSize)

	for i := 0; i < UBufDataSize; i++ {
		writeData[i] = byte(i)
		readData[i] = 0
	}

	_, err := ub.Read(readData)
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	n, err := ub.Write(writeData)
	if err == nil {
		t.Fatal("Unexpected write result")
	}

	if n != UBufDataSize {
		t.Fatal("Unexpected write result")
	}

	if ub.TailWritableLength() != 0 {
		t.Fatal("Unexpected tail writable length")
	}

	if ub.ReadableLength() != UBufDataSize {
		t.Fatal("Unexpected tail readable length")
	}

	nread, err := ub.Read(readData)
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if nread != UBufDataSize {
		t.Fatal("Unexpected read result")
	}

	if ub.TailWritableLength() != 0 {
		t.Fatal("Unexpected tail writable length")
	}

	if ub.ReadableLength() != 0 {
		t.Fatal("Unexpected tail readable length")
	}

	for i := 0; i < UBufDataSize; i++ {
		if readData[i] != byte(i) {
			t.Fatal("Unexpected read data")
		}
	}

	ub.Reset()

	_, err = ub.Write(writeData)
	if err == nil {
		t.Fatal("Unexpected write result")
	}

	_, err = ub.Write(writeData)
	if err == nil {
		t.Fatal("Unexpected write result")
	}

	nread, err = ub.Read(make([]byte, 128))
	if err == nil {
		t.Fatal("Unexpected read result")
	}

	if nread != UBufDataSize {
		t.Fatal("Unexpected read result")
	}

	ub.Reset()

	if ub.TailWritableLength() != UBufDataSize {
		t.Fatal("Unexpected tail writable length")
	}

	if ub.ReadableLength() != 0 {
		t.Fatal("Unexpected readable length")
	}
}

type DummyReaderWriter int

func (drw *DummyReaderWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("DummyReaderWriter")
}

func (drw *DummyReaderWriter) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("DummyReaderWriter")
}

func TestReadFrom(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufReserved)
	var drw DummyReaderWriter = 0

	_, err := ub.ReadFrom(&drw)
	if err == nil {
		t.Fatal("Unexpected read result")
	}

	writeData := make([]byte, UBufDataSize)
	readData := make([]byte, UBufDataSize)

	for i := 0; i < UBufDataSize; i++ {
		writeData[i] = byte(i)
		readData[i] = 0
	}

	br := bytes.NewReader(writeData)

	n, err := ub.ReadFrom(br)
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if n != UBufDataSize {
		t.Fatal("Unexpected read result")
	}

	if ub.TailWritableLength() != 0 {
		t.Fatal("Unexpected tail writable length")
	}

	if ub.ReadableLength() != UBufDataSize {
		t.Fatal("Unexpected readable length")
	}

	_, err = ub.ReadFrom(br)
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	nread, err := ub.Read(readData)
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if nread != UBufDataSize {
		t.Fatal("Unexpected read result")
	}

	for i := 0; i < UBufDataSize; i++ {
		if readData[i] != byte(i) {
			t.Fatal("Unexpected read data")
		}
	}
}

func TestWriteTo(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufReserved)

	var drw DummyReaderWriter = 0

	_, err := ub.WriteTo(&drw)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	writeData := make([]byte, UBufDataSize)
	for i := 0; i < UBufDataSize; i++ {
		writeData[i] = byte(i)
	}

	bb := bytes.NewBuffer(make([]byte, 0))

	n, err := ub.Write(writeData)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	if n != UBufDataSize {
		t.Fatal("Unexpected write result")
	}

	_, err = ub.WriteTo(&drw)
	if err == nil {
		t.Fatal("Unexpected write result")
	}

	nwrite, err := ub.WriteTo(bb)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	if nwrite != UBufDataSize {
		t.Fatal("Unexpected write result")
	}

	data := bb.Bytes()
	for i := 0; i < UBufDataSize; i++ {
		if data[i] != byte(i) {
			t.Fatal("Unexpected read data")
		}
	}
}

func TestReadWriteUxx(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufReserved)

	err := ub.WriteByte(1)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	err = ub.WriteU16(0x2345)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	err = ub.WriteU32(0x6789abcd)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	err = ub.WriteU64(0x1357246875318642)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	b, err := ub.ReadByte()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if b != 1 {
		t.Fatal("Unexpected read data")
	}

	u16, err := ub.ReadU16()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u16 != 0x2345 {
		t.Fatal("Unexpected read data")
	}

	u32, err := ub.ReadU32()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u32 != 0x6789abcd {
		t.Fatal("Unexpected read data")
	}

	u64, err := ub.ReadU64()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u64 != 0x1357246875318642 {
		t.Fatal("Unexpected read data")
	}

	ub.Reset()

	if _, err := ub.ReadU16(); err == nil {
		t.Fatal("Unexpected read result")
	}

	if _, err := ub.ReadU32(); err == nil {
		t.Fatal("Unexpected read result")
	}

	if _, err := ub.ReadU64(); err == nil {
		t.Fatal("Unexpected read result")
	}

	if _, err := ub.ReadU16BE(); err == nil {
		t.Fatal("Unexpected read result")
	}

	if _, err := ub.ReadU32BE(); err == nil {
		t.Fatal("Unexpected read result")
	}

	if _, err := ub.ReadU64BE(); err == nil {
		t.Fatal("Unexpected read result")
	}

	if _, err := ub.ReadU16LE(); err == nil {
		t.Fatal("Unexpected read result")
	}

	if _, err := ub.ReadU32LE(); err == nil {
		t.Fatal("Unexpected read result")
	}

	if _, err := ub.ReadU64LE(); err == nil {
		t.Fatal("Unexpected read result")
	}

	ub = UBufAllocWithHeadReserved(UBufCapacity, UBufCapacity)

	if err := ub.WriteU16(1); err == nil {
		t.Fatal("Unexpected Write result")
	}

	if err := ub.WriteU32(1); err == nil {
		t.Fatal("Unexpected Write result")
	}

	if err := ub.WriteU64(1); err == nil {
		t.Fatal("Unexpected Write result")
	}

	if err := ub.WriteU16BE(1); err == nil {
		t.Fatal("Unexpected Write result")
	}

	if err := ub.WriteU32BE(1); err == nil {
		t.Fatal("Unexpected Write result")
	}

	if err := ub.WriteU64BE(1); err == nil {
		t.Fatal("Unexpected Write result")
	}

	if err := ub.WriteU16LE(1); err == nil {
		t.Fatal("Unexpected Write result")
	}

	if err := ub.WriteU32LE(1); err == nil {
		t.Fatal("Unexpected Write result")
	}

	if err := ub.WriteU64LE(1); err == nil {
		t.Fatal("Unexpected Write result")
	}
}

func TestReadWriteUxxBE(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufReserved)

	err := ub.WriteU16BE(0x2345)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	err = ub.WriteU32BE(0x6789abcd)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	err = ub.WriteU64BE(0x1357246875318642)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	u16, err := ub.ReadU16BE()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u16 != 0x2345 {
		t.Fatal("Unexpected read data")
	}

	u32, err := ub.ReadU32BE()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u32 != 0x6789abcd {
		t.Fatal("Unexpected read data")
	}

	u64, err := ub.ReadU64BE()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u64 != 0x1357246875318642 {
		t.Fatal("Unexpected read data")
	}
}

func TestReadWriteUxxLE(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufReserved)

	err := ub.WriteU16LE(0x2345)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	err = ub.WriteU32LE(0x6789abcd)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	err = ub.WriteU64LE(0x1357246875318642)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	u16, err := ub.ReadU16LE()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u16 != 0x2345 {
		t.Fatal("Unexpected read data")
	}

	u32, err := ub.ReadU32LE()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u32 != 0x6789abcd {
		t.Fatal("Unexpected read data")
	}

	u64, err := ub.ReadU64LE()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u64 != 0x1357246875318642 {
		t.Fatal("Unexpected read data")
	}
}

func TestReadWriteHeadUxxNoSpace(t *testing.T) {
	ub := UBufAlloc(UBufCapacity)

	err := ub.WriteHeadByte(1)
	if err == nil {
		t.Fatal("Unexpected WriteHead result")
	}

	if err := ub.WriteHeadU16(1); err == nil {
		t.Fatal("Unexpected WriteHead result")
	}

	if err := ub.WriteHeadU32(1); err == nil {
		t.Fatal("Unexpected WriteHead result")
	}

	if err := ub.WriteHeadU64(1); err == nil {
		t.Fatal("Unexpected WriteHead result")
	}

	if err := ub.WriteHeadU16BE(1); err == nil {
		t.Fatal("Unexpected WriteHead result")
	}

	if err := ub.WriteHeadU32BE(1); err == nil {
		t.Fatal("Unexpected WriteHead result")
	}

	if err := ub.WriteHeadU64BE(1); err == nil {
		t.Fatal("Unexpected WriteHead result")
	}

	if err := ub.WriteHeadU16LE(1); err == nil {
		t.Fatal("Unexpected WriteHead result")
	}

	if err := ub.WriteHeadU32LE(1); err == nil {
		t.Fatal("Unexpected WriteHead result")
	}

	if err := ub.WriteHeadU64LE(1); err == nil {
		t.Fatal("Unexpected WriteHead result")
	}
}

func TestReadWriteHeadBytes(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufCapacity)

	err := ub.WriteHeadBytes([]byte{1, 2})
	if err != nil {
		t.Fatal("Unexpected WriteHead result")
	}

	bytes := make([]byte, 2)

	n, err := ub.Read(bytes)

	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if n != 2 {
		t.Fatal("Unexpected read result")
	}

	if bytes[0] != 1 && bytes[1] != 2 {
		t.Fatal("Unexpected read result")
	}
}

func TestReadWriteHeadUxx(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufCapacity)

	err := ub.WriteHeadByte(1)
	if err != nil {
		t.Fatal("Unexpected WriteHead result")
	}

	err = ub.WriteHeadU16(0x2345)
	if err != nil {
		t.Fatal("Unexpected WriteHead result")
	}

	err = ub.WriteHeadU32(0x6789abcd)
	if err != nil {
		t.Fatal("Unexpected WriteHead result")
	}

	err = ub.WriteHeadU64(0x1357246875318642)
	if err != nil {
		t.Fatal("Unexpected WriteHead result")
	}

	u64, err := ub.ReadU64()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u64 != 0x1357246875318642 {
		t.Fatal("Unexpected read data")
	}

	u32, err := ub.ReadU32()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u32 != 0x6789abcd {
		t.Fatal("Unexpected read data")
	}

	u16, err := ub.ReadU16()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u16 != 0x2345 {
		t.Fatal("Unexpected read data")
	}

	b, err := ub.ReadByte()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if b != 1 {
		t.Fatal("Unexpected read data")
	}
}

func TestReadWriteHeadUxxBE(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufCapacity)

	err := ub.WriteHeadU16BE(0x2345)
	if err != nil {
		t.Fatal("Unexpected WriteHead result")
	}

	err = ub.WriteHeadU32BE(0x6789abcd)
	if err != nil {
		t.Fatal("Unexpected WriteHead result")
	}

	err = ub.WriteHeadU64BE(0x1357246875318642)
	if err != nil {
		t.Fatal("Unexpected WriteHead result")
	}

	u64, err := ub.ReadU64BE()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u64 != 0x1357246875318642 {
		t.Fatal("Unexpected read data")
	}

	u32, err := ub.ReadU32BE()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u32 != 0x6789abcd {
		t.Fatal("Unexpected read data")
	}

	u16, err := ub.ReadU16BE()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u16 != 0x2345 {
		t.Fatal("Unexpected read data")
	}
}

func TestReadWriteHeadUxxLE(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufCapacity)

	err := ub.WriteHeadU16LE(0x2345)
	if err != nil {
		t.Fatal("Unexpected WriteHead result")
	}

	err = ub.WriteHeadU32LE(0x6789abcd)
	if err != nil {
		t.Fatal("Unexpected WriteHead result")
	}

	err = ub.WriteHeadU64LE(0x1357246875318642)
	if err != nil {
		t.Fatal("Unexpected WriteHead result")
	}

	u64, err := ub.ReadU64LE()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u64 != 0x1357246875318642 {
		t.Fatal("Unexpected read data")
	}

	u32, err := ub.ReadU32LE()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u32 != 0x6789abcd {
		t.Fatal("Unexpected read data")
	}

	u16, err := ub.ReadU16LE()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u16 != 0x2345 {
		t.Fatal("Unexpected read data")
	}
}

func TestPeekBytes(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufReserved)
	bytes := make([]byte, 2)

	_, err := ub.PeekByte()
	if err == nil {
		t.Fatal("Unexpected peek result")
	}

	_, err = ub.Peek(bytes)
	if err == nil {
		t.Fatal("Unexpected peek result")
	}

	_, err = ub.Write([]byte{1, 2, 3})
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	_, err = ub.Peek(make([]byte, 4))
	if err == nil {
		t.Fatal("Unexpected peek result")
	}

	b, err := ub.PeekByte()
	if err != nil {
		t.Fatal("Unexpected peek result")
	}

	if b != 1 {
		t.Fatal("Unexpected peek data")
	}

	n, err := ub.Peek(bytes)
	if err != nil {
		t.Fatal("Unexpected peek result")
	}

	if n != 2 {
		t.Fatal("Unexpected peek data")
	}

	if bytes[0] != 1 || bytes[1] != 2 {
		t.Fatal("Unexpected peek data")
	}
}

func TestPeekUxx(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufReserved)

	_, err := ub.PeekU16()
	if err == nil {
		t.Fatal("Unexpected peek result")
	}

	err = ub.WriteU16(0x2345)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	u16, err := ub.PeekU16()
	if err != nil {
		t.Fatal("Unexpected peek result")
	}

	if u16 != 0x2345 {
		t.Fatal("Unexpected peek data")
	}

	ub.Reset()

	_, err = ub.PeekU32()
	if err == nil {
		t.Fatal("Unexpected peek result")
	}

	err = ub.WriteU32(0x6789abcd)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	u32, err := ub.PeekU32()
	if err != nil {
		t.Fatal("Unexpected peek result")
	}

	if u32 != 0x6789abcd {
		t.Fatal("Unexpected peek data")
	}

	ub.Reset()

	_, err = ub.PeekU64()
	if err == nil {
		t.Fatal("Unexpected peek result")
	}

	err = ub.WriteU64(0x1357246875318642)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	u64, err := ub.PeekU64()
	if err != nil {
		t.Fatal("Unexpected peek result")
	}

	if u64 != 0x1357246875318642 {
		t.Fatal("Unexpected peek data")
	}
}

func TestPeekUxxBE(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufReserved)

	_, err := ub.PeekU16BE()
	if err == nil {
		t.Fatal("Unexpected peek result")
	}

	err = ub.WriteU16BE(0x2345)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	u16, err := ub.PeekU16BE()
	if err != nil {
		t.Fatal("Unexpected peek result")
	}

	if u16 != 0x2345 {
		t.Fatal("Unexpected peek data")
	}

	ub.Reset()

	_, err = ub.PeekU32BE()
	if err == nil {
		t.Fatal("Unexpected peek result")
	}

	err = ub.WriteU32BE(0x6789abcd)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	u32, err := ub.PeekU32BE()
	if err != nil {
		t.Fatal("Unexpected peek result")
	}

	if u32 != 0x6789abcd {
		t.Fatal("Unexpected peek data")
	}

	ub.Reset()

	_, err = ub.PeekU64BE()
	if err == nil {
		t.Fatal("Unexpected peek result")
	}

	err = ub.WriteU64BE(0x1357246875318642)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	u64, err := ub.PeekU64BE()
	if err != nil {
		t.Fatal("Unexpected peek result")
	}

	if u64 != 0x1357246875318642 {
		t.Fatal("Unexpected peek data")
	}
}

func TestPeekUxxLE(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufReserved)

	_, err := ub.PeekU16LE()
	if err == nil {
		t.Fatal("Unexpected peek result")
	}

	err = ub.WriteU16LE(0x2345)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	u16, err := ub.PeekU16LE()
	if err != nil {
		t.Fatal("Unexpected peek result")
	}

	if u16 != 0x2345 {
		t.Fatal("Unexpected peek data")
	}

	ub.Reset()

	_, err = ub.PeekU32LE()
	if err == nil {
		t.Fatal("Unexpected peek result")
	}

	err = ub.WriteU32LE(0x6789abcd)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	u32, err := ub.PeekU32LE()
	if err != nil {
		t.Fatal("Unexpected peek result")
	}

	if u32 != 0x6789abcd {
		t.Fatal("Unexpected peek data")
	}

	ub.Reset()

	_, err = ub.PeekU64LE()
	if err == nil {
		t.Fatal("Unexpected peek result")
	}

	err = ub.WriteU64LE(0x1357246875318642)
	if err != nil {
		t.Fatal("Unexpected write result")
	}

	u64, err := ub.PeekU64LE()
	if err != nil {
		t.Fatal("Unexpected peek result")
	}

	if u64 != 0x1357246875318642 {
		t.Fatal("Unexpected peek data")
	}
}

func TestSnapshotReference(t *testing.T) {
	if UBufMakeSnapshot(nil, 999) != nil {
		t.Fatal("Unexpected returned snap")
	}

	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufReserved)
	ub.WriteU16(0x1234)

	if UBufMakeSnapshot(ub, 999) != nil {
		t.Fatal("Unexpected returned snap")
	}

	snap := UBufMakeSnapshot(ub, UBufSnapshotTypeReference)

	u16, err := snap.ReadU16()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u16 != 0x1234 {
		t.Fatal("Unexpected read data")
	}

	ub.WriteU16(0x5678)

	_, err = snap.ReadU16()
	if err == nil {
		t.Fatal("Unexpected read result")
	}
}

func TestSnapshotCopyOnWrite(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufReserved)
	ub.WriteU16(0x1234)

	snap := UBufMakeSnapshot(ub, UBufSnapshotTypeCopyOnWrite)

	if &snap.data.bytes != &ub.data.bytes {
		t.Fatal("did copy")
	}

	snap.WriteU16(0x5678)

	ub.WriteU16(0x8765)

	u16, err := snap.ReadU16()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u16 != 0x1234 {
		t.Fatal("Unexpected read data")
	}

	u16, err = snap.ReadU16()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u16 != 0x5678 {
		t.Fatal("Unexpected read data")
	}

	u16, err = ub.ReadU16()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u16 != 0x1234 {
		t.Fatal("Unexpected read data")
	}

	u16, err = ub.ReadU16()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u16 != 0x8765 {
		t.Fatal("Unexpected read data")
	}

	if &snap.data.bytes == &ub.data.bytes {
		t.Fatal("did not copy on write")
	}
}

func TestSnapshotCopyDriectly(t *testing.T) {
	ub := UBufAllocWithHeadReserved(UBufCapacity, UBufReserved)
	ub.WriteU16(0x1234)

	snap := UBufMakeSnapshot(ub, UBufSnapshotTypeCopyDirectly)

	if &snap.data.bytes == &ub.data.bytes {
		t.Fatal("did not copy directly")
	}

	snap.WriteU16(0x5678)

	ub.WriteU16(0x8765)

	u16, err := snap.ReadU16()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u16 != 0x1234 {
		t.Fatal("Unexpected read data")
	}

	u16, err = snap.ReadU16()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u16 != 0x5678 {
		t.Fatal("Unexpected read data")
	}

	u16, err = ub.ReadU16()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u16 != 0x1234 {
		t.Fatal("Unexpected read data")
	}

	u16, err = ub.ReadU16()
	if err != nil {
		t.Fatal("Unexpected read result")
	}

	if u16 != 0x8765 {
		t.Fatal("Unexpected read data")
	}

	if &snap.data.bytes == &ub.data.bytes {
		t.Fatal("did not copy on write")
	}
}
