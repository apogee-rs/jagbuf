package jagbuf

import (
	"io"
)

type Buffer struct {
	data       []byte
	readIndex  int
	writeIndex int
}

// NewBuffer creates a new buffer with an initial capacity of 64.
func NewBuffer() *Buffer {
	return NewWithCapacity(64)
}

// NewWithCapacity creates a new buffer with the provided
// initial capacity.
func NewWithCapacity(capacity int) *Buffer {
	return &Buffer{
		data:       make([]byte, capacity),
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
	// For zero value buffers
	if b.data == nil {
		b.data = make([]byte, numBytes)
		return
	}

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
// Start (inclusive) to End (exclusive) - [start .. end)
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

// Write will copy data into the buffer, and grow the underlying buffer
// if there is not enough space.
func (b *Buffer) Write(data []byte) int {
	b.ensureWritable(len(data))

	written := copy(b.data[b.writeIndex:], data)

	defer func() { b.writeIndex += written }()
	return written
}
