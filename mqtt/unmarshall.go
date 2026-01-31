package mqtt

import (
	"encoding/binary"
	"io"
)

type Unmarshall struct {
	buffer io.Reader
	err    error
}

func NewUnmarshall(r io.Reader) *Unmarshall {
	return &Unmarshall{buffer: r}
}

func (u *Unmarshall) Error() error {
	return u.err
}

func (u *Unmarshall) ReadBytes(n int) []byte {
	if u.err != nil {
		return []byte{}
	}

	buf := make([]byte, n)
	_, u.err = u.buffer.Read(buf)
	return buf
}

func (u *Unmarshall) ReadByte() byte {
	if u.err != nil {
		return 0
	}

	val := make([]byte, 1)
	_, err := u.buffer.Read(val)
	if err != nil {
		u.err = err
	}

	return val[0]
}

func (u *Unmarshall) ReadUint16() uint16 {
	if u.err != nil {
		return 0
	}

	buf := u.ReadBytes(2)
	if u.err != nil {
		return 0
	}

	return binary.BigEndian.Uint16(buf)
}

func (u *Unmarshall) ReadString() string {
	if u.err != nil {
		return ""
	}

	strLength := u.ReadUint16()
	if u.err != nil {
		return ""
	}

	buf := u.ReadBytes(int(strLength))
	if u.err != nil {
		return ""
	}

	return string(buf)
}
