package jagbuf

import "io"

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

func (b *Buffer) ReadUint64LE() (uint64, error) {
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

func (b *Buffer) ReadInt64LE() (int64, error) {
	val, err := b.ReadUint64LE()
	return int64(val), err
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

func (b *Buffer) WriteInt64(v int64) {
	b.WriteUint64(uint64(v))
}

func (b *Buffer) WriteUint64LE(v uint64) {
	b.ensureWritable(8)

	b.data[b.writeIndex] = byte(v & 0xFF)
	b.data[b.writeIndex+1] = byte(v >> 8)
	b.data[b.writeIndex+2] = byte(v >> 16)
	b.data[b.writeIndex+3] = byte(v >> 24)
	b.data[b.writeIndex+4] = byte(v >> 32)
	b.data[b.writeIndex+5] = byte(v >> 40)
	b.data[b.writeIndex+6] = byte(v >> 48)
	b.data[b.writeIndex+7] = byte(v >> 56)

	defer func() { b.writeIndex += 8 }()
}

func (b *Buffer) WriteInt64LE(v int64) {
	b.WriteUint64LE(uint64(v))
}
