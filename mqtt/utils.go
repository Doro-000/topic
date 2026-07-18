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

func PacketFactory(header MqttHeader) GenericPacket {
	switch header.PacketType {
	case CONNECT:
		return &MqttConnect{
			MqttHeader: header,
		}
	case CONNACK:
		return &MqttConnectAck{
			MqttHeader: header,
		}
	case PUBLISH:
		return &MqttPublish{
			MqttHeader: header,
		}
	case PUBACK, PUBREC, PUBREL, PUBCOMP, UNSUBACK:
		return &MqttSimplePacket{
			MqttHeader: header,
		}
	case SUBSCRIBE:
		return &MqttSubscribe{
			MqttHeader: header,
		}
	case SUBACK:
		return &MqttSubAck{
			MqttHeader: header,
		}
	case UNSUBSCRIBE:
		return &MqttUnsubscribe{
			MqttHeader: header,
		}
	case DISCONNECT, PINGREQ, PINGRESP:
		return &MqttHeaderOnlyPacket{
			MqttHeader: header,
		}

	case RESERVED_PACKET_TYPE, RESERVED_PACKET_TYPE_2:
	default:
		return nil
	}

	return nil
}
