package mqtt

import "io"

type MqttSubAckCode byte

const (
	SUCCESS_MAX_QOS_0 MqttSubAckCode = iota
	SUCCESS_MAX_QOS_1
	SUCCESS_MAX_QOS_2
	FAILURE = 0x80
)

type MqttSubAck struct {
	MqttHeader
	PacketIdentifier uint16
	Payload          []MqttSubAckCode
}

func (ack *MqttSubAck) Marshall(marshaller *Marshall) error {
	ack.MqttHeader.Marshall(marshaller)

	remainingLen := EncodeRemainingLength(2 + len(ack.Payload))
	marshaller.WriteBytes(remainingLen)
	marshaller.WriteUint16(ack.PacketIdentifier)

	for _, code := range ack.Payload {
		marshaller.WriteByte(byte(code))
	}

	return marshaller.Error()
}

func (ack *MqttSubAck) Unmarshall(unmarshaller *Unmarshall) error {
	ack.PacketIdentifier = unmarshaller.ReadUint16()

	codes, _ := io.ReadAll(unmarshaller.buffer)

	for _, code := range codes {
		ack.Payload = append(ack.Payload, MqttSubAckCode(code))
	}

	return unmarshaller.Error()
}

func (packet *MqttSubAck) GetType() MQTTControlPacketType {
	return SUBACK
}
