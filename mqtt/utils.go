package mqtt

import (
	"errors"
	"io"
)

func DecodeRemainingLen(packet io.Reader) (int, error) {
	multiplier := 1
	value := 0

	b := make([]byte, 1)
	for {
		_, err := packet.Read(b)

		if err != nil {
			return 0, err
		}

		nextByte := int(b[0])
		value += (nextByte & 127) * multiplier

		if (nextByte & 128) == 0 {
			break
		}

		if multiplier > MAX_MULTIPLIER_REMAIN_LEN {
			return 0, errors.New("Malformed Remaining Length")
		}

		multiplier *= 128
	}

	return value, nil
}

func EncodeRemainingLength(len int) (encodedBytes []byte) {
	if len == 0 {
		return []byte{0}
	}

	for len > 0 {
		encodedByte := byte(len % 128)
		len /= 128

		if len > 0 {
			encodedByte |= 0x08
		}

		encodedBytes = append(encodedBytes, encodedByte)
	}
	return
}
