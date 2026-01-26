package mqtt

import (
	"encoding/binary"
	"io"
)

type Marshall struct {
	buffer io.Writer
	err    error
}

func NewMarshall(w io.Writer) *Marshall {
	return &Marshall{buffer: w}
}

func (e *Marshall) Error() error {
	return e.err
}

func (e *Marshall) WriteBytes(bytes []byte) {
	if e.err != nil {
		return
	}

	_, e.err = e.buffer.Write(bytes)
}

func (e *Marshall) WriteByte(val byte) error {
	if e.err != nil {
		return e.err
	}

	_, e.err = e.buffer.Write([]byte{val})
	return e.err
}

func (e *Marshall) WriteUint16(val uint16) {
	if e.err != nil {
		return
	}
	buf := make([]byte, 0, 2)
	e.WriteBytes(binary.BigEndian.AppendUint16(buf, val))
}

func (e *Marshall) WriteString(str string) {
	if e.err != nil {
		return
	}

	strLen := uint16(len(str))
	e.WriteUint16(strLen)

	if e.err != nil {
		return
	}

	e.WriteBytes([]byte(str))
}
