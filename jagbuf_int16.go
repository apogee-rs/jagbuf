package jagbuf

import "io"

func (b *Buffer) ReadUint16() (uint16, error) {
	if b.ReadableBytes() < 2 {
		return 0, io.EOF
	}

	val := uint16(b.data[b.readIndex]) << 8
	val |= uint16(b.data[b.readIndex+1])

	defer func() { b.readIndex += 2 }()
	return val, nil
}

// ReadUint16_Sub reads an uint16 from the buffer and applies the `value - 128` transform to the
// lower order bits.
func (b *Buffer) ReadUint16_Sub() (uint16, error) {
	if b.ReadableBytes() < 2 {
		return 0, io.EOF
	}

	val := uint16(b.data[b.readIndex]) << 8
	val |= uint16(b.data[b.readIndex+1]) - 128

	defer func() { b.readIndex += 2 }()
	return val, nil
}

func (b *Buffer) ReadUint16LE() (uint16, error) {
	if b.ReadableBytes() < 2 {
		return 0, io.EOF
	}

	val := uint16(b.data[b.readIndex])
	val |= uint16(b.data[b.readIndex+1]) << 8

	defer func() { b.readIndex += 2 }()
	return val, nil
}

// ReadUint16LE_Sub reads an uint16 from the buffer and applies the `value - 128` transform to the
// higher order bits.
func (b *Buffer) ReadUint16LE_Sub() (uint16, error) {
	if b.ReadableBytes() < 2 {
		return 0, io.EOF
	}

	val := uint16(b.data[b.readIndex]) - 128
	val |= uint16(b.data[b.readIndex+1]) << 8

	defer func() { b.readIndex += 2 }()
	return val, nil
}

func (b *Buffer) ReadInt16() (int16, error) {
	val, err := b.ReadUint16()
	return int16(val), err
}

func (b *Buffer) ReadInt16LE() (int16, error) {
	val, err := b.ReadUint16LE()
	return int16(val), err
}

func (b *Buffer) WriteUint16(v uint16) {
	b.ensureWritable(2)

	b.data[b.writeIndex] = byte(v >> 8)
	b.data[b.writeIndex+1] = byte(v & 0xFF)

	defer func() { b.writeIndex += 2 }()
}

func (b *Buffer) WriteInt16(v int16) {
	b.WriteUint16(uint16(v))
}

func (b *Buffer) WriteUint16LE(v uint16) {
	b.ensureWritable(2)

	b.data[b.writeIndex] = byte(v & 0xFF)
	b.data[b.writeIndex+1] = byte(v >> 8)

	defer func() { b.writeIndex += 2 }()
}

func (b *Buffer) WriteInt16LE(v int16) {
	b.WriteUint16LE(uint16(v))
}
