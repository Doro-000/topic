package mqtt

type MqttPublishAck struct {
	Header           MqttHeader
	PacketIdentifier uint16
}

func (ack *MqttPublishAck) Marshall(marshaller *Marshall) error {
	marshaller.WriteUint16(ack.PacketIdentifier)

	return marshaller.Error()
}

func (ack *MqttPublishAck) UnmarshallMqttPublishAck(unmarshaller *Unmarshall) error {
	ack.PacketIdentifier = unmarshaller.ReadUint16()

	return unmarshaller.Error()
}
