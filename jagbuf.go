package jagbuf

import "io"

type Buffer struct {
	data       []byte
	readIndex  int
	writeIndex int
}

// NewBuffer creates a new buffer with an initial capacity of 64.
func NewBuffer() *Buffer {
	return NewBufferWithCapacity(64)
}

// NewBufferWithCapacity creates a new buffer with the provided
// initial capacity.
func NewBufferWithCapacity(capacity int) *Buffer {
	return &Buffer{
		data:       make([]byte, 0, capacity),
		readIndex:  0,
		writeIndex: 0,
	}
}

// Wrap creates a new buffer copying the provided data.
// This sets the write index to the end of the data, and
// the read index to the beginning of the data. This does
// not take ownership of the provided data.
func Wrap(data []byte) *Buffer {
	return &Buffer{
		data:       append(make([]byte, 0, len(data)), data...),
		readIndex:  0,
		writeIndex: len(data),
	}
}

func (b *Buffer) Capacity() int {
	return cap(b.data)
}

// Grow allocates a new buffer with capacity expanded by n and copies
// the existing data into it.
func (b *Buffer) Grow(n int) {
	b.data = append(make([]byte, len(b.data), cap(b.data)+n), b.data...)
}

func (b *Buffer) ensureWritable(numBytes int) {
	for b.WritableBytes() < numBytes {
		// Double our capacity until we have enough room for numBytes
		b.Grow(b.Capacity())
	}
}

func (b *Buffer) ReadableBytes() int {
	return b.writeIndex - b.readIndex
}

func (b *Buffer) WritableBytes() int {
	return b.Capacity() - b.writeIndex
}

// Bytes returns a new slice with only the data written (0 to writeIndex).
// This frees the underlying data to be reused or modified.
func (b *Buffer) Bytes() []byte {
	return append(make([]byte, 0, b.writeIndex), b.data[:b.writeIndex]...)
}

// Slice returns a new slice with a copy of the requested data.
func (b *Buffer) Slice(start int, end int) []byte {
	return append(make([]byte, 0, end-start), b.data[start:end]...)
}

// Reset returns both the read and write indexes to zero, allowing any
// existing data to be re-read or overwritten. Use ResetReadIndex or
// ResetWriteIndex to reset one or the other.
func (b *Buffer) Reset() {
	b.ResetWriteIndex()
	b.ResetReadIndex()
}

// ResetReadIndex returns the read index to zero, allowing any existing
// data to be re-read.
func (b *Buffer) ResetReadIndex() {
	b.readIndex = 0
}

// ResetWriteIndex returns the write index to zero, allowing any existing
// data to be overwritten.
func (b *Buffer) ResetWriteIndex() {
	b.writeIndex = 0
}

func (b *Buffer) Skip(n int) {
	b.readIndex += n
}

// Read will copy data from the buffer and place it into dst. This will
// return the number of bytes read.
func (b *Buffer) Read(dst []byte) (int, error) {
	if b.ReadableBytes() < 1 {
		return 0, io.EOF
	}

	read := copy(dst, b.data[b.readIndex:])

	defer func() { b.readIndex += read }()
	return read, nil
}

func (b *Buffer) ReadUint8() (uint8, error) {
	if b.ReadableBytes() < 1 {
		return 0, io.EOF
	}

	val := b.data[b.readIndex]

	defer func() { b.readIndex += 1 }()
	return val, nil
}

func (b *Buffer) ReadInt8() (int8, error) {
	val, err := b.ReadUint8()
	return int8(val), err
}

func (b *Buffer) ReadUint16() (uint16, error) {
	if b.ReadableBytes() < 2 {
		return 0, io.EOF
	}

	val := uint16(b.data[b.readIndex]) << 8
	val |= uint16(b.data[b.readIndex+1])

	defer func() { b.readIndex += 2 }()
	return val, nil
}

func (b *Buffer) ReadLEUint16() (uint16, error) {
	if b.ReadableBytes() < 2 {
		return 0, io.EOF
	}

	val := uint16(b.data[b.readIndex])
	val |= uint16(b.data[b.readIndex+1]) << 8

	defer func() { b.readIndex += 2 }()
	return val, nil
}

func (b *Buffer) ReadInt16() (int16, error) {
	val, err := b.ReadUint16()
	return int16(val), err
}

func (b *Buffer) ReadLEInt16() (int16, error) {
	val, err := b.ReadLEUint16()
	return int16(val), err
}

func (b *Buffer) ReadUint24() (uint32, error) {
	if b.ReadableBytes() < 3 {
		return 0, io.EOF
	}

	val := uint32(b.data[b.readIndex]) << 16
	val |= uint32(b.data[b.readIndex+1]) << 8
	val |= uint32(b.data[b.readIndex+2])

	defer func() { b.readIndex += 3 }()
	return val, nil
}

func (b *Buffer) ReadLEUint24() (uint32, error) {
	if b.ReadableBytes() < 3 {
		return 0, io.EOF
	}

	val := uint32(b.data[b.readIndex])
	val |= uint32(b.data[b.readIndex+1]) << 8
	val |= uint32(b.data[b.readIndex+2]) << 16

	defer func() { b.readIndex += 3 }()
	return val, nil
}

func (b *Buffer) ReadInt24() (int32, error) {
	val, err := b.ReadUint24()
	return int32(val), err
}

func (b *Buffer) ReadLEInt24() (int32, error) {
	if b.ReadableBytes() < 3 {
		return 0, io.EOF
	}

	// TODO: Can we just cast ReadLEUint24 to an int32?
	// Or will overflows not be handled correctly?
	val := int32(b.data[b.readIndex])
	val |= int32(b.data[b.readIndex+1]) << 8
	val |= int32(b.data[b.readIndex+2]) << 16

	defer func() { b.readIndex += 3 }()
	return val, nil
}

func (b *Buffer) ReadUint32() (uint32, error) {
	if b.ReadableBytes() < 4 {
		return 0, io.EOF
	}

	val := uint32(b.data[b.readIndex]) << 24
	val |= uint32(b.data[b.readIndex+1]) << 16
	val |= uint32(b.data[b.readIndex+2]) << 8
	val |= uint32(b.data[b.readIndex+3])

	defer func() { b.readIndex += 4 }()
	return val, nil
}

func (b *Buffer) ReadLEUint32() (uint32, error) {
	if b.ReadableBytes() < 4 {
		return 0, io.EOF
	}

	val := uint32(b.data[b.readIndex])
	val |= uint32(b.data[b.readIndex+1]) << 8
	val |= uint32(b.data[b.readIndex+2]) << 16
	val |= uint32(b.data[b.readIndex+3]) << 24

	defer func() { b.readIndex += 4 }()
	return val, nil
}

func (b *Buffer) ReadInt32() (int32, error) {
	val, err := b.ReadUint32()
	return int32(val), err
}

func (b *Buffer) ReadLEInt32() (int32, error) {
	val, err := b.ReadLEUint32()
	return int32(val), err
}

func (b *Buffer) ReadUint64() (uint64, error) {
	if b.ReadableBytes() < 8 {
		return 0, io.EOF
	}

	val := uint64(b.data[b.readIndex]) << 56
	val |= uint64(b.data[b.readIndex+1]) << 48
	val |= uint64(b.data[b.readIndex+2]) << 40
	val |= uint64(b.data[b.readIndex+3]) << 32
	val |= uint64(b.data[b.readIndex+4]) << 24
	val |= uint64(b.data[b.readIndex+5]) << 16
	val |= uint64(b.data[b.readIndex+6]) << 8
	val |= uint64(b.data[b.readIndex+7])

	defer func() { b.readIndex += 8 }()
	return val, nil
}

func (b *Buffer) ReadLEUint64() (uint64, error) {
	if b.ReadableBytes() < 8 {
		return 0, io.EOF
	}

	val := uint64(b.data[b.readIndex])
	val |= uint64(b.data[b.readIndex+1]) << 8
	val |= uint64(b.data[b.readIndex+2]) << 16
	val |= uint64(b.data[b.readIndex+3]) << 24
	val |= uint64(b.data[b.readIndex+4]) << 32
	val |= uint64(b.data[b.readIndex+5]) << 40
	val |= uint64(b.data[b.readIndex+6]) << 48
	val |= uint64(b.data[b.readIndex+7]) << 56

	defer func() { b.readIndex += 8 }()
	return val, nil
}

func (b *Buffer) ReadInt64() (int64, error) {
	val, err := b.ReadUint64()
	return int64(val), err
}

func (b *Buffer) ReadLEInt64() (int64, error) {
	val, err := b.ReadLEUint64()
	return int64(val), err
}

// Write will copy data into the buffer, and grow the underlying buffer
// if there is not enough space.
func (b *Buffer) Write(data []byte) int {
	if b.WritableBytes() < len(data) {
		// Grow by however many bytes we are short
		b.Grow(len(data) - b.WritableBytes())
	}

	written := copy(b.data[b.writeIndex:], data)

	defer func() { b.writeIndex += written }()
	return written
}

func (b *Buffer) WriteUint8(v uint8) {
	b.ensureWritable(1)

	b.data[b.writeIndex] = v

	defer func() { b.writeIndex += 1 }()
}

func (b *Buffer) WriteUint16(v uint16) {
	b.ensureWritable(2)

	b.data[b.writeIndex] = byte(v >> 8)
	b.data[b.writeIndex+1] = byte(v & 0xFF)

	defer func() { b.writeIndex += 2 }()
}

func (b *Buffer) WriteUint24(v uint32) {
	b.ensureWritable(3)

	b.data[b.writeIndex] = byte(v >> 16)
	b.data[b.writeIndex+1] = byte(v >> 8)
	b.data[b.writeIndex+2] = byte(v & 0xFF)

	defer func() { b.writeIndex += 3 }()
}

func (b *Buffer) WriteUint32(v uint32) {
	b.ensureWritable(4)

	b.data[b.writeIndex] = byte(v >> 24)
	b.data[b.writeIndex+1] = byte(v >> 16)
	b.data[b.writeIndex+2] = byte(v >> 8)
	b.data[b.writeIndex+3] = byte(v & 0xFF)

	defer func() { b.writeIndex += 4 }()
}

func (b *Buffer) WriteUint64(v uint64) {
	b.ensureWritable(8)

	b.data[b.writeIndex] = byte(v >> 56)
	b.data[b.writeIndex+1] = byte(v >> 48)
	b.data[b.writeIndex+2] = byte(v >> 40)
	b.data[b.writeIndex+3] = byte(v >> 32)
	b.data[b.writeIndex+4] = byte(v >> 24)
	b.data[b.writeIndex+5] = byte(v >> 16)
	b.data[b.writeIndex+6] = byte(v >> 8)
	b.data[b.writeIndex+7] = byte(v & 0xFF)

	defer func() { b.writeIndex += 8 }()
}
