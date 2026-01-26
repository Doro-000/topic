package mqtt

import "io"

type MqttPublish struct {
	Header           MqttHeader
	TopicName        string
	PacketIdentifier uint16
	Payload          []byte
}

func (packet *MqttPublish) Marshall(marshaller *Marshall) error {
	marshaller.WriteString(packet.TopicName)

	if packet.Header.Qos == AT_LEAST_ONCE || packet.Header.Qos == EXACTLY_ONCE {
		marshaller.WriteUint16(packet.PacketIdentifier)
	}

	marshaller.WriteBytes(packet.Payload)

	return marshaller.Error()
}

func (packet *MqttPublish) Unmarshall(unmarshaller *Unmarshall) error {
	packet.TopicName, _ = unmarshaller.ReadString()

	if packet.Header.Qos == AT_LEAST_ONCE || packet.Header.Qos == EXACTLY_ONCE {
		packet.PacketIdentifier = unmarshaller.ReadUint16()
	}

	// The rest is payload
	// We need to know the total length of the variable header + payload
	// This should be handled by the main decode loop, which should use an io.LimitedReader
	// For now, we assume the unmarshaller's buffer is already limited.
	payload, err := io.ReadAll(unmarshaller.buffer)
	if err != nil {
		unmarshaller.err = err
	}
	packet.Payload = payload

	return unmarshaller.Error()
}
