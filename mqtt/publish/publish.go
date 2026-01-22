package publish

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	BaseMqtt "github.com/Doro-000/topic/mqtt"
)

type MqttPublish struct {
	Header           BaseMqtt.MqttHeader
	TopicName        string
	PacketIdentifier []byte
	Payload          []byte
}

// implement generic packet functions
func (packet *MqttPublish) GetHeader() byte {
	return packet.Header.Value
}

func (packet *MqttPublish) GetVariableHeader() []byte {
	topicNameLen := []byte{}
	topicNameLen = binary.BigEndian.AppendUint16(topicNameLen, uint16(len(packet.TopicName)))

	res := []byte{}
	res = append(res, topicNameLen...)
	res = append(res, []byte(packet.TopicName)...)
	res = append(res, packet.PacketIdentifier...)

	return res
}

func (packet *MqttPublish) GetPayload() []byte {
	return packet.Payload
}

// implement unmarshalling function
func UnmarshallMqttPublish(header BaseMqtt.MqttHeader, packet io.Reader) (*MqttPublish, error) {
	if header.GetType() != BaseMqtt.PUBLISH {
		return nil, fmt.Errorf("Called Connect unmarshall on packet type: %v", header.GetType())
	}

	res := MqttPublish{Header: header}

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

	topicName, payloadLen := packetUnmarshall.String()
	res.TopicName = topicName

	qosLevel := header.GetQos()
	if qosLevel == BaseMqtt.AT_LEAST_ONCE || qosLevel == BaseMqtt.EXACTLY_ONCE {
		res.PacketIdentifier = []byte{packetUnmarshall.Uint8(), packetUnmarshall.Uint8()}
		payloadLen += 2
	}

	payload := make([]byte, len-payloadLen)
	_, err = io.ReadFull(varHeaderAndPayload, payload)
	if err != nil {
		return nil, err
	}

	res.Payload = payload

	return &res, nil
}
