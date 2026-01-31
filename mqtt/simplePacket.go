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
