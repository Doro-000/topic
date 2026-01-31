package mqtt

import "io"

type MqttPublish struct {
	MqttHeader
	TopicName        string
	PacketIdentifier uint16
	Payload          []byte
}

func (packet *MqttPublish) Marshall(marshaller *Marshall) error {
	/** Header */
	packet.MqttHeader.Marshall(marshaller)

	topicNameLen := (2 + len(packet.TopicName))
	packetIDlen := 2
	remainingLength := EncodeRemainingLength(packetIDlen + topicNameLen + len(packet.Payload))
	marshaller.WriteBytes(remainingLength)

	/** Variable Header */
	marshaller.WriteString(packet.TopicName)

	if packet.Qos == AT_LEAST_ONCE || packet.Qos == EXACTLY_ONCE {
		marshaller.WriteUint16(packet.PacketIdentifier)
	}

	/** Payload */
	marshaller.WriteBytes(packet.Payload)

	return marshaller.Error()
}

func (packet *MqttPublish) Unmarshall(unmarshaller *Unmarshall) error {
	packet.TopicName = unmarshaller.ReadString()

	if packet.Qos == AT_LEAST_ONCE || packet.Qos == EXACTLY_ONCE {
		packet.PacketIdentifier = unmarshaller.ReadUint16()
	}

	payload, err := io.ReadAll(unmarshaller.buffer)
	if err != nil {
		unmarshaller.err = err
	}
	packet.Payload = payload

	return unmarshaller.Error()
}
