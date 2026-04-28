package mqtt

import (
	"fmt"
	"io"
)

type MqttUnsubscribe struct {
	MqttHeader
	PacketIdentifier uint16
	Payload          []string
}

func (packet *MqttUnsubscribe) getPayloadLen() int {
	payloadLen := 0
	for _, topicFilter := range packet.Payload {
		payloadLen += 2 + len(topicFilter)
	}

	return payloadLen
}

func (packet *MqttUnsubscribe) Marshall(marshaller *Marshall) error {
	/** Header */
	packet.MqttHeader.Marshall(marshaller)

	packetIDlen := 2
	remainingLength := EncodeRemainingLength(packetIDlen + packet.getPayloadLen())
	marshaller.WriteBytes(remainingLength)

	/** Variable Header */
	marshaller.WriteUint16(packet.PacketIdentifier)

	for _, topicFilter := range packet.Payload {
		marshaller.WriteString(topicFilter)
	}

	return marshaller.Error()
}

func (packet *MqttUnsubscribe) Unmarshall(unmarshaller *Unmarshall) error {
	packet.PacketIdentifier = unmarshaller.ReadUint16()
	packet.Payload = make([]string, 0)

	for {
		topicFilter := unmarshaller.ReadString()

		if unmarshaller.Error() == io.EOF {
			break
		}

		packet.Payload = append(packet.Payload, topicFilter)
	}

	if len(packet.Payload) == 0 {
		return fmt.Errorf("Unsubscribe with no payload found !")
	}

	if unmarshaller.Error() != io.EOF {
		return unmarshaller.Error()
	}
	return nil
}

func (packet *MqttUnsubscribe) GetType() MQTTControlPacketType {
	return UNSUBSCRIBE
}
