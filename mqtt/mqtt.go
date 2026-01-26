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
	PUBLISH
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
	RESERVED_PACKET_TYPE_2
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

type GenericPacket interface {
	Marshall(e *Marshall) error
	Unmarshall(u *Unmarshall) error
}

type MqttHeader struct {
	Retain     bool
	Qos        QoSLevel
	Dup        bool
	PacketType MQTTControlPacketType
}

func NewMqttHeader(packetType MQTTControlPacketType, dup bool, qos QoSLevel, retain bool) (*MqttHeader, error) {
	if qos == RESERVED_QOS_LEVEL {
		return nil, errors.New("Bad QoS level")
	}

	if packetType == RESERVED_PACKET_TYPE || packetType == RESERVED_PACKET_TYPE_2 {
		return nil, fmt.Errorf("Packet Type %v not allowed", packetType)
	}

	return &MqttHeader{
		Retain:     retain,
		Qos:        qos,
		Dup:        dup,
		PacketType: packetType,
	}, nil
}

func (h *MqttHeader) MarshallMqttHeader() (byte, error) {
	if h.PacketType == RESERVED_PACKET_TYPE || h.PacketType == RESERVED_PACKET_TYPE_2 {
		return 0, fmt.Errorf("Packet Type %v not allowed", h.PacketType)
	}

	if h.Qos == RESERVED_QOS_LEVEL {
		return 0, errors.New("Bad QoS level in header")
	}

	var val byte = 0

	if h.Retain {
		val |= 1
	}

	if h.Qos != AT_LEAST_ONCE {
		val |= byte(h.Qos) << 1
	}

	if h.Dup {
		val |= (1 << 3)
	}

	val |= byte(h.PacketType) << 4

	return val, nil
}

func UnmarshallMqttHeader(value byte) (*MqttHeader, error) {
	packetType := MQTTControlPacketType(value >> 4)
	qosLevel := QoSLevel((value >> 1) & 3)
	var retain bool = false
	var dup bool = false

	if value&1 == 1 {
		retain = true
	}

	if ((value >> 3) & 1) == 1 {
		dup = true
	}

	return NewMqttHeader(
		packetType,
		dup,
		qosLevel,
		retain,
	)
}
