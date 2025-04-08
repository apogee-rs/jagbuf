package jagbuf

import "io"

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

// ReadUint8_Sub reads an uint8 from the buffer and applies the `value - 128` transform.
func (b *Buffer) ReadUint8_Sub() (uint8, error) {
	if b.ReadableBytes() < 1 {
		return 0, io.EOF
	}

	val := b.data[b.readIndex] - 128

	defer func() { b.readIndex += 1 }()
	return val, nil
}

// ReadUint8_Neg reads an uint8 from the buffer and applies the `0 - value` transform.
func (b *Buffer) ReadUint8_Neg() (uint8, error) {
	if b.ReadableBytes() < 1 {
		return 0, io.EOF
	}

	val := 0 - b.data[b.readIndex]

	defer func() { b.readIndex += 1 }()
	return val, nil
}

// ReadUint8_Mirror reads an uint8 from the buffer and applies the `128 - value` transform.
func (b *Buffer) ReadUint8_Mirror() (uint8, error) {
	if b.ReadableBytes() < 1 {
		return 0, io.EOF
	}

	val := 128 - b.data[b.readIndex]

	defer func() { b.readIndex += 1 }()
	return val, nil
}

// ReadInt8_Sub reads an int8 from the buffer and applies the `value - 128` transform.
func (b *Buffer) ReadInt8_Sub() (int8, error) {
	val, err := b.ReadUint8_Sub()
	return int8(val), err
}

// ReadInt8_Neg reads an int8 from the buffer and applies the `0 - value` transform.
func (b *Buffer) ReadInt8_Neg() (int8, error) {
	val, err := b.ReadUint8_Neg()
	return int8(val), err
}

// ReadInt8_Mirror reads an int8 from the buffer and applies the `128 - value` transform.
func (b *Buffer) ReadInt8_Mirror() (int8, error) {
	val, err := b.ReadUint8_Mirror()
	return int8(val), err
}

func (b *Buffer) WriteUint8(v uint8) {
	b.ensureWritable(1)

	b.data[b.writeIndex] = v

	defer func() { b.writeIndex += 1 }()
}

func (b *Buffer) WriteInt8(v int8) {
	b.WriteUint8(uint8(v))
}
