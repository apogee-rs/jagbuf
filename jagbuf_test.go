package jagbuf

import (
	"bytes"
	"fmt"
	"testing"
)

func TestWrap(t *testing.T) {
	data := []byte("test")
	buffer := Wrap(data)

	data = []byte("other")
	if bytes.Equal(data, buffer.Bytes()) {
		t.Error("Wrap fail: underlying data was changed")
	}
}

func TestNewBufferWithCapacity(t *testing.T) {
	buffer := NewWithCapacity(128)

	if buffer.Capacity() != 128 {
		t.Error("NewWithCapacity fail: underlying buffer capacity is wrong")
	}
}

func TestBuffer_Grow(t *testing.T) {
	buffer := NewWithCapacity(64)
	buffer.Grow(64)

	if buffer.Capacity() != 128 {
		t.Error("NewWithCapacity fail: underlying buffer capacity is wrong")
	}
}

func TestBuffer_Write(t *testing.T) {
	buffer := NewWithCapacity(64)

	data := []byte{0x0, 0x1, 0x2, 0x3}
	_ = buffer.Write(data)

	if !bytes.Equal(buffer.Bytes(), data) {
		t.Error("Write fail: underlying buffer mismatch")
	}
}

func TestBuffer_Slice(t *testing.T) {
	buffer := Wrap([]byte{0x0, 0x1, 0x2, 0x3})

	slice := buffer.Slice(1, 3)

	buffer.writeIndex = 1
	buffer.WriteUint8(0x0)
	buffer.WriteUint8(0x0)

	if bytes.Equal(slice, []byte{0x0, 0x0}) {
		t.Error("Slice fail: slice was modified")
	}
}

func TestBuffer_ReadString(t *testing.T) {
	data := append([]byte("hello world"), 0x0)
	buffer := Wrap(data)

	str, err := buffer.ReadString()
	if err != nil {
		t.Fatal(err)
	}

	if str != "hello world" {
		t.Errorf("ReadString fail: Expected \"%s\" but received \"%s\"", "hello world", str)
	}
}

func TestBuffer_ReadJagString(t *testing.T) {
	data := append([]byte{0x0}, []byte("hello world")...)
	data = append(data, 0x0)
	buffer := Wrap(data)

	str, err := buffer.ReadJagString()
	if err != nil {
		t.Fatal(err)
	}

	if str != "hello world" {
		t.Errorf("ReadJagString fail: Expected \"%s\" but received \"%s\"", "hello world", str)
	}
}

func TestBuffer_ReadWriteIntegrity(t *testing.T) {
	buffer := NewWithCapacity(64)

	buffer.WriteUint32(0x10203040)

	val, err := buffer.ReadUint32()
	if err != nil {
		t.Fatal(err)
	}

	if val != 0x10203040 {
		t.Errorf("Read/Write fail: Expected 0x10203040 but received 0x%x", val)
	}
}

func TestBuffer_ReadWriteMixedEndianness(t *testing.T) {
	buffer := NewWithCapacity(64)

	buffer.WriteUint32(0x10203040)

	val, err := buffer.ReadUint32LE()
	if err != nil {
		t.Fatal(err)
	}

	if val != 0x40302010 {
		t.Errorf("Read/Write fail: Expected 0x40302010 but received 0x%x", val)
	}
}

func TestBuffer_ReadInt24(t *testing.T) {
	buffer := NewWithCapacity(64)

	buffer.WriteInt24(8388608) //overflows 24-bit

	val, err := buffer.ReadInt24()
	if err != nil {
		t.Fatal(err)
	}

	if val != -8388608 {
		fmt.Printf("ReadInt24 fail: Expected -8388608 but received %d", val)
	}
}

func TestBuffer_ReadInt24_WithNeg(t *testing.T) {
	buffer := NewWithCapacity(64)

	buffer.WriteInt24(-8388609) //under flows 24-bit

	val, err := buffer.ReadInt24()
	if err != nil {
		t.Fatal(err)
	}

	if val != 8388607 {
		fmt.Printf("ReadInt24 fail: Expected 8388607 but received %d", val)
	}
}

func TestBuffer_ReadInt24LE(t *testing.T) {
	buffer := NewWithCapacity(64)

	buffer.WriteInt24LE(8388608) //overflows 24-bit

	val, err := buffer.ReadInt24LE()
	if err != nil {
		t.Fatal(err)
	}

	if val != -8388608 {
		fmt.Printf("ReadInt24LE fail: Expected -8388608 but received %d", val)
	}
}

func TestBuffer_ReadInt24LE_WithNeg(t *testing.T) {
	buffer := NewWithCapacity(64)

	buffer.WriteInt24LE(-8388609) //under flows 24-bit

	val, err := buffer.ReadInt24LE()
	if err != nil {
		t.Fatal(err)
	}

	if val != 8388607 {
		fmt.Printf("ReadInt24LE fail: Expected 8388607 but received %d (0x%x)", val, val)
	}
}

func TestBuffer_WriteInt32(t *testing.T) {
	buffer := NewWithCapacity(64)

	buffer.WriteInt32(-32768)

	val, err := buffer.ReadInt32()
	if err != nil {
		t.Fatal(err)
	}

	if val != -32768 {
		fmt.Printf("WriteInt32 fail: Expected -32768 but received %d", val)
	}
}
