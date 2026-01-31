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
	PacketIdentifier uint16
	payload          []MqttSubAckCode
}

func (ack *MqttSubAck) Marshall(marshaller *Marshall) error {
	marshaller.WriteUint16(ack.PacketIdentifier)

	for _, code := range ack.payload {
		marshaller.WriteByte(byte(code))
	}

	return marshaller.Error()
}

func (ack *MqttSubAck) Unmarshall(unmarshaller *Unmarshall) error {
	ack.PacketIdentifier = unmarshaller.ReadUint16()

	codes, _ := io.ReadAll(unmarshaller.buffer)

	for _, code := range codes {
		ack.payload = append(ack.payload, MqttSubAckCode(code))
	}

	return unmarshaller.Error()
}
