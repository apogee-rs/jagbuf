package jagbuf

import "io"

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

func (b *Buffer) ReadUint32LE() (uint32, error) {
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

// ReadUint32V1 reads an int32 from the buffer with a special Jagex endianness.
// This is equivalent to big endian but with the first 2 bytes shifted to the end.
func (b *Buffer) ReadUint32V1() (uint32, error) {
	if b.ReadableBytes() < 4 {
		return 0, io.EOF
	}

	val := uint32(b.data[b.readIndex]) << 8
	val |= uint32(b.data[b.readIndex+1])
	val |= uint32(b.data[b.readIndex+2]) << 24
	val |= uint32(b.data[b.readIndex+3]) << 16

	defer func() { b.readIndex += 4 }()
	return val, nil
}

// ReadUint32V2 reads an int32 from the buffer with a special Jagex endianness.
// This is equivalent to little endian but with the first 2 bytes shifted to the end.
func (b *Buffer) ReadUint32V2() (uint32, error) {
	if b.ReadableBytes() < 4 {
		return 0, io.EOF
	}

	val := uint32(b.data[b.readIndex]) << 16
	val |= uint32(b.data[b.readIndex+1]) << 24
	val |= uint32(b.data[b.readIndex+2])
	val |= uint32(b.data[b.readIndex+3]) << 8

	defer func() { b.readIndex += 4 }()
	return val, nil
}

func (b *Buffer) ReadInt32() (int32, error) {
	val, err := b.ReadUint32()
	return int32(val), err
}

func (b *Buffer) ReadInt32LE() (int32, error) {
	val, err := b.ReadUint32LE()
	return int32(val), err
}

// ReadInt32V1 reads an int32 from the buffer with a special Jagex endianness.
// This is equivalent to big endian but with the first 2 bytes shifted to the end.
func (b *Buffer) ReadInt32V1() (int32, error) {
	val, err := b.ReadUint32V1()
	return int32(val), err
}

// ReadInt32V2 reads an int32 from the buffer with a special Jagex endianness.
// This is equivalent to little endian but with the first 2 bytes shifted to the end.
func (b *Buffer) ReadInt32V2() (int32, error) {
	val, err := b.ReadUint32V2()
	return int32(val), err
}

func (b *Buffer) WriteUint32(v uint32) {
	b.ensureWritable(4)

	b.data[b.writeIndex] = byte(v >> 24)
	b.data[b.writeIndex+1] = byte(v >> 16)
	b.data[b.writeIndex+2] = byte(v >> 8)
	b.data[b.writeIndex+3] = byte(v & 0xFF)

	defer func() { b.writeIndex += 4 }()
}

func (b *Buffer) WriteInt32(v int32) {
	b.WriteUint32(uint32(v))
}

func (b *Buffer) WriteUint32LE(v uint32) {
	b.ensureWritable(4)

	b.data[b.writeIndex] = byte(v & 0xFF)
	b.data[b.writeIndex+1] = byte(v >> 8)
	b.data[b.writeIndex+2] = byte(v >> 16)
	b.data[b.writeIndex+3] = byte(v >> 24)

	defer func() { b.writeIndex += 4 }()
}

func (b *Buffer) WriteInt32LE(v int32) {
	b.WriteUint32LE(uint32(v))
}

// WriteUint32V1 writes an uint32 to the buffer using a special Jagex
// endianness. This is equivalent to big endian, with the first 2 bytes
// shuffled to the end.
func (b *Buffer) WriteUint32V1(v uint32) {
	b.ensureWritable(4)

	b.data[b.writeIndex] = byte(v >> 8)
	b.data[b.writeIndex+1] = byte(v & 0xFF)
	b.data[b.writeIndex+2] = byte(v >> 24)
	b.data[b.writeIndex+3] = byte(v >> 16)

	defer func() { b.writeIndex += 4 }()
}

// WriteInt32V1 writes an int32 to the buffer using a special Jagex
// endianness. This is equivalent to big endian, with the first 2 bytes
// shuffled to the end.
func (b *Buffer) WriteInt32V1(v int32) {
	b.WriteUint32V1(uint32(v))
}

// WriteUint32V2 writes an uint32 to the buffer using a special Jagex
// endianness. This is equivalent to little endian, with the first 2 bytes
// shuffled to the end.
func (b *Buffer) WriteUint32V2(v uint32) {
	b.ensureWritable(4)

	b.data[b.writeIndex] = byte(v >> 16)
	b.data[b.writeIndex+1] = byte(v >> 24)
	b.data[b.writeIndex+2] = byte(v & 0xFF)
	b.data[b.writeIndex+3] = byte(v >> 8)

	defer func() { b.writeIndex += 4 }()
}

// WriteInt32V2 writes an int32 to the buffer using a special Jagex
// endianness. This is equivalent to little endian, with the first 2 bytes
// shuffled to the end.
func (b *Buffer) WriteInt32V2(v int32) {
	b.WriteUint32V2(uint32(v))
}
