package mqtt

import (
	"fmt"
)

type MqttConnAckCode byte

const (
	ACCEPTED MqttConnAckCode = iota
	BAD_PROTO_V
	ID_REJECTED
	SERVER_DED
	BAD_UNAME_PASS
	NOT_AUTHORIZED
)

type MqttConnectAck struct {
	Header         MqttHeader
	SessionPresent bool
	ReturnCode     MqttConnAckCode
}

func (ack *MqttConnectAck) Marshall(marshaller *Marshall) error {
	var sp byte
	if ack.SessionPresent {
		sp = 1
	}
	marshaller.WriteByte(sp)
	marshaller.WriteByte(byte(ack.ReturnCode))

	return marshaller.Error()
}

func (ack *MqttConnectAck) Unmarshall(unmarshaller *Unmarshall) error {
	sp, _ := unmarshaller.ReadByte()
	if unmarshaller.Error() != nil {
		return unmarshaller.Error()
	}

	// Validate that reserved bits (1-7) are 0
	if (sp & 0xFE) != 0 {
		return fmt.Errorf("malformed connack packet: reserved bits of acknowledge flags must be 0")
	}

	ack.SessionPresent = (sp & 0x01) == 1

	returnCode, _ := unmarshaller.ReadByte()
	if unmarshaller.Error() != nil {
		return unmarshaller.Error()
	}
	ack.ReturnCode = MqttConnAckCode(returnCode)

	// Validate the return code
	if ack.ReturnCode > NOT_AUTHORIZED {
		return fmt.Errorf("malformed connack packet: invalid return code %d", ack.ReturnCode)
	}

	return nil
}
