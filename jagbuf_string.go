package jagbuf

import (
	"errors"
	"io"
	"strings"
)

func (b *Buffer) ReadString() (string, error) {
	if b.ReadableBytes() < 1 {
		return "", io.EOF
	}

	builder := &strings.Builder{}
	for i := 0; b.data[b.readIndex+i] != 0; i++ {
		builder.WriteByte(b.data[b.readIndex+i])
	}

	defer func() { b.readIndex += builder.Len() + 1 }()
	return builder.String(), nil
}

func (b *Buffer) ReadJagString() (string, error) {
	if b.ReadableBytes() < 1 {
		return "", io.EOF
	}

	peek := b.data[b.readIndex]
	b.readIndex++

	if peek != 0 {
		return "", errors.New("jagstring read: expected byte to be 0 in position 0")
	}

	return b.ReadString()
}
