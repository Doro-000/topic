package mqtt

import (
	"fmt"
	"io"
	"maps"
)

type MqttSubscribe struct {
	MqttHeader
	PacketIdentifier uint16
	Payload          map[string]QoSLevel
}

func (packet *MqttSubscribe) getPayloadLen() int {
	payloadLen := 0
	for topicFilter := range maps.Keys(packet.Payload) {
		payloadLen += 1 // 1 byte for the qos level
		payloadLen += 2 + len(topicFilter)
	}

	return payloadLen
}

func (packet *MqttSubscribe) Marshall(marshaller *Marshall) error {
	/** Header */
	packet.MqttHeader.Marshall(marshaller)

	packetIDlen := 2
	remainingLength := EncodeRemainingLength(packetIDlen + packet.getPayloadLen())
	marshaller.WriteBytes(remainingLength)

	/** Variable Header */
	marshaller.WriteUint16(packet.PacketIdentifier)

	for topicFilter, qos := range packet.Payload {
		marshaller.WriteString(topicFilter)
		marshaller.WriteByte(byte(qos))
	}

	return marshaller.Error()
}

func (packet *MqttSubscribe) Unmarshall(unmarshaller *Unmarshall) error {
	packet.PacketIdentifier = unmarshaller.ReadUint16()

	packet.Payload = make(map[string]QoSLevel)

	for {
		if unmarshaller.Error() == io.EOF {
			break
		}

		topicFilter := unmarshaller.ReadString()
		topicQos := unmarshaller.ReadByte()

		packet.Payload[topicFilter] = QoSLevel(topicQos)
	}

	if len(packet.Payload) == 0 {
		return fmt.Errorf("Subscribe with no payload found !")
	}

	return unmarshaller.Error()
}

func (packet *MqttSubscribe) GetType() MQTTControlPacketType {
	return SUBSCRIBE
}
