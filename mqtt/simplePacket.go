package mqtt

/*
Represents packets with no payload and just a packet identifier:
  - PUBACK
  - PUBREC
  - PUBREL
  - PUBCOMP
  - UNSUBACK
*/
type MqttSimplePacket struct {
	MqttHeader
	PacketIdentifier uint16
}

func (p *MqttSimplePacket) Marshall(marshaller *Marshall) error {
	p.MqttHeader.Marshall(marshaller)
	marshaller.WriteByte(0x02)

	marshaller.WriteUint16(p.PacketIdentifier)

	return marshaller.Error()
}

func (p *MqttSimplePacket) Unmarshall(unmarshaller *Unmarshall) error {
	p.PacketIdentifier = unmarshaller.ReadUint16()

	return unmarshaller.Error()
}

func (packet *MqttSimplePacket) GetType() MQTTControlPacketType {
	return packet.MqttHeader.PacketType
}

/*
Represents packets with just a header:
  - DISCONNECT
  - PINGREQ
  - PINGRESP
*/
type MqttHeaderOnlyPacket struct {
	MqttHeader
}

func (p *MqttHeaderOnlyPacket) Marshall(marshaller *Marshall) error {
	p.MqttHeader.Marshall(marshaller)
	marshaller.WriteByte(0x00)
	return marshaller.Error()
}

func (p *MqttHeaderOnlyPacket) Unmarshall(unmarshaller *Unmarshall) error {
	return nil
}

func (packet *MqttHeaderOnlyPacket) GetType() MQTTControlPacketType {
	return packet.MqttHeader.PacketType
}
