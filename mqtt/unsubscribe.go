package mqtt

import "io"

type MqttUnsubscribe struct {
	MqttHeader
	PacketIdentifier uint16
	Payload          []string
}

func (packet *MqttUnsubscribe) Marshall(marshaller *Marshall) error {
	/** Header */
	packet.MqttHeader.Marshall(marshaller)

	payloadLen := len(packet.Payload) + (2 * len(packet.Payload))
	packetIDlen := 2
	remainingLength := EncodeRemainingLength(packetIDlen + payloadLen)
	marshaller.WriteBytes(remainingLength)

	marshaller.WriteUint16(packet.PacketIdentifier)

	for _, topicFilter := range packet.Payload {
		marshaller.WriteString(topicFilter)
	}

	return marshaller.Error()
}

func (packet *MqttUnsubscribe) Unmarshall(unmarshaller *Unmarshall) error {
	packet.PacketIdentifier = unmarshaller.ReadUint16()

	for {
		if unmarshaller.Error() == io.EOF {
			break
		}

		topicFilter := unmarshaller.ReadString()
		packet.Payload = append(packet.Payload, topicFilter)
	}

	return unmarshaller.Error()
}
