package jagbuf

import "io"

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

func (b *Buffer) ReadUint24LE() (uint32, error) {
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
	if b.ReadableBytes() < 3 {
		return 0, io.EOF
	}

	var val int32 = 0
	if (b.data[b.readIndex] & 0x80) != 0 {
		//sign extend if negative
		val = -8388608
	}

	val |= int32(b.data[b.readIndex]) << 16
	val |= int32(b.data[b.readIndex+1]) << 8
	val |= int32(b.data[b.readIndex+2])

	defer func() { b.readIndex += 3 }()
	return val, nil
}

func (b *Buffer) ReadInt24LE() (int32, error) {
	if b.ReadableBytes() < 3 {
		return 0, io.EOF
	}

	var val int32 = 0
	if (b.data[b.readIndex+2] & 0x80) != 0 {
		//sign extend if negative
		val = -8388608
	}

	val |= int32(b.data[b.readIndex+2]) << 16
	val |= int32(b.data[b.readIndex+1]) << 8
	val |= int32(b.data[b.readIndex])

	defer func() { b.readIndex += 3 }()
	return val, nil
}

func (b *Buffer) WriteUint24(v uint32) {
	b.ensureWritable(3)

	b.data[b.writeIndex] = byte(v >> 16)
	b.data[b.writeIndex+1] = byte(v >> 8)
	b.data[b.writeIndex+2] = byte(v & 0xFF)

	defer func() { b.writeIndex += 3 }()
}

func (b *Buffer) WriteInt24(v int32) {
	b.WriteUint24(uint32(v))
}

func (b *Buffer) WriteUint24LE(v uint32) {
	b.ensureWritable(3)

	b.data[b.writeIndex] = byte(v & 0xFF)
	b.data[b.writeIndex+1] = byte(v >> 8)
	b.data[b.writeIndex+2] = byte(v >> 16)

	defer func() { b.writeIndex += 3 }()
}

func (b *Buffer) WriteInt24LE(v int32) {
	b.WriteUint24LE(uint32(v))
}
