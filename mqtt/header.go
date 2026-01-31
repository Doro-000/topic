package mqtt

import (
	"errors"
	"fmt"
)

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

func (h *MqttHeader) Marshall(marshaller *Marshall) error {
	var val byte = 0

	if h.Retain {
		val |= 1
	}

	if h.Qos != AT_MOST_ONCE {
		val |= byte(h.Qos) << 1
	}

	if h.Dup {
		val |= (1 << 3)
	}

	val |= byte(h.PacketType) << 4

	marshaller.WriteByte(val)
	return marshaller.Error()
}

func UnmarshallMqttHeader(value byte) (*MqttHeader, error) {
	packetType := MQTTControlPacketType(value >> 4)

	dup := false
	if ((value >> 3) & 1) == 1 {
		dup = true
	}

	qos := QoSLevel((value >> 1) & 3)

	retain := false
	if value&1 == 1 {
		retain = true
	}

	return NewMqttHeader(packetType, dup, qos, retain)
}
