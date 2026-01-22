package publishack

import (
	"bytes"
	"fmt"
	"io"

	BaseMqtt "github.com/Doro-000/topic/mqtt"
)

type MqttPublishAck struct {
	Header           BaseMqtt.MqttHeader
	PacketIdentifier []byte
}

func (packet *MqttPublishAck) GetHeader() byte {
	return packet.Header.Value
}

func (packet *MqttPublishAck) GetVariableHeader() []byte {
	return packet.PacketIdentifier
}

func (packet *MqttPublishAck) GetPayload() []byte {
	return []byte{}
}

func UnmarshallMqttPublishAck(header BaseMqtt.MqttHeader, packet io.Reader) (*MqttPublishAck, error) {
	if header.GetType() != BaseMqtt.PUBACK {
		return nil, fmt.Errorf("Called Connect unmarshall on packet type: %v", header.GetType())
	}

	res := MqttPublishAck{Header: header}

	// decode remaining length
	len, err := BaseMqtt.DecodeRemainingLen(packet)

	if err != nil {
		return nil, err
	}

	remainingPacket := make([]byte, len)
	_, err = io.ReadFull(packet, remainingPacket)

	if err != nil {
		return nil, err
	}

	varHeaderAndPayload := bytes.NewReader(remainingPacket)
	packetUnmarshall := BaseMqtt.NewUnmarshall(varHeaderAndPayload)

	res.PacketIdentifier = []byte{packetUnmarshall.Uint8(), packetUnmarshall.Uint8()}

	return &res, nil
}
