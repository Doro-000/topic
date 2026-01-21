package mqtt

import (
	"errors"
	"fmt"
)

type MQTTControlPacketType byte
type QoSLevel byte

// Constants

//go:generate stringer -type=MQTTControlPacketType
const (
	RESERVED_PACKET_TYPE MQTTControlPacketType = iota
	CONNECT
	CONNACK
	PUBACK
	PUBREC  // Publish received
	PUBREL  // Publish release
	PUBCOMP // Publish complete
	SUBSCRIBE
	SUBACK
	UNSUBSCRIBE
	UNSUBACK
	PINGREQ
	PINGRESP
	DISCONNECT
)

const (
	AT_MOST_ONCE QoSLevel = iota
	AT_LEAST_ONCE
	EXACTLY_ONCE
	RESERVED_QOS_LEVEL
)

const MQTT_MAX_PACKET_SIZE = 1500 // 64kB
const MQTT_REMAIN_LEN_MAX = 4     // maximum of 4 bytes for Remaining Length
const MAX_MULTIPLIER_REMAIN_LEN = 128 * 128 * 128

type MqttHeader struct {
	Value byte
}

type GenericPacket interface {
	GetHeader() byte
	GetVariableHeader() []byte
	GetPayload() []byte
}

func UnmarshallMqttHeader(value byte) (*MqttHeader, error) {
	packetType := MQTTControlPacketType(value >> 4)

	if packetType == RESERVED_PACKET_TYPE {
		return nil, fmt.Errorf("Packet Type %v not allowed", packetType)
	}

	if (value>>1)&3 >= 3 {
		return nil, errors.New("Bad QoS level in header")
	}

	return &MqttHeader{Value: value}, nil
}

func NewMqttHeader(packetType MQTTControlPacketType, retain bool, qos QoSLevel, dup bool) (*MqttHeader, error) {
	var value byte = 0

	if retain {
		value |= 1
	}

	if qos == RESERVED_QOS_LEVEL {
		return nil, errors.New("Bad QoS level")
	} else {
		value |= (byte(qos) << 1)
	}

	if dup {
		value |= 8
	}

	if packetType == RESERVED_PACKET_TYPE {
		return nil, fmt.Errorf("Packet Type %v not allowed", packetType)
	}

	value |= byte(packetType) << 4

	return &MqttHeader{Value: value}, nil
}

// 0th bit
func (h *MqttHeader) GetRetain() bool {
	return (h.Value & 0x01) == 1
}

// 1st and 2nd bits
func (h *MqttHeader) GetQos() byte {
	return h.Value & 0x03
}

// 3rd bit
func (h *MqttHeader) GetDup() byte {
	return (h.Value >> 3) & 0x01
}

// 4th - 7th bits
func (h *MqttHeader) GetType() MQTTControlPacketType {
	nibble := (h.Value >> 4) & 0x0F

	if nibble == 15 {
		return RESERVED_PACKET_TYPE
	}
	return MQTTControlPacketType(nibble)
}
