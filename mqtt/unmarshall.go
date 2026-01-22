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

func (packet *Unmarshall) Error() error {
	return packet.err
}

func (packet *Unmarshall) Uint8() uint8 {
	if packet.err != nil {
		return 0
	}

	res := make([]byte, 1)
	_, packet.err = io.ReadFull(packet.buffer, res)
	return res[0]
}

func (packet *Unmarshall) Uint16() uint16 {
	if packet.err != nil {
		return 0
	}

	res := make([]byte, 2)
	_, packet.err = io.ReadFull(packet.buffer, res)
	return binary.BigEndian.Uint16(res)
}

func (packet *Unmarshall) String() (string, int) {
	if packet.err != nil {
		return "", 0
	}

	strLength := packet.Uint16()
	if packet.err != nil {
		return "", 0
	}

	str := make([]byte, strLength)
	_, packet.err = io.ReadFull(packet.buffer, str)
	return string(str), int(strLength) + 2 // 2 bytes for the len
}
