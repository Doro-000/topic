package connectack

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	BaseMqtt "github.com/Doro-000/topic/mqtt"
)

type MqttConnAckCodes byte

const (
	ACCEPTED MqttConnAckCodes = iota
	BAD_PROTO_V
	ID_REJECTED
	SERVER_DED
	BAD_UNAME_PASS
	NOT_AUTHORIZED
)

type MqttConnectAck struct {
	Header         BaseMqtt.MqttHeader
	SessionPresent bool
	ReturnCode     MqttConnAckCodes
}

func (packet *MqttConnectAck) GetHeader() byte {
	return packet.Header.Value
}

func (packet *MqttConnectAck) GetVariableHeader() []byte {
	var flag byte

	if packet.SessionPresent {
		flag = 1
	} else {
		flag = 0
	}

	return []byte{flag, byte(packet.ReturnCode)}
}

func (packet *MqttConnectAck) GetPayload() []byte {
	return []byte{}
}

func UnmarshallMqttConnectAck(header BaseMqtt.MqttHeader, packet io.Reader) (*MqttConnectAck, error) {
	if header.GetType() != BaseMqtt.CONNACK {
		return nil, fmt.Errorf("Called ConnectAck unmarshall on packet type: %v", header.GetType())
	}

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

	sessionPresentVal := packetUnmarshall.Uint8()
	returnCode := MqttConnAckCodes(packetUnmarshall.Uint8())

	if packetUnmarshall.Error() != nil {
		return nil, packetUnmarshall.Error()
	}

	if sessionPresentVal > 1 {
		return nil, errors.New("Invalid connect acknowledge flag")
	}

	var sessionPresent bool = false
	if sessionPresentVal == 1 {
		sessionPresent = true
	}

	return &MqttConnectAck{
		Header:         header,
		SessionPresent: sessionPresent,
		ReturnCode:     returnCode,
	}, nil
}

func NewMqttConnectAck(header BaseMqtt.MqttHeader, sessionPresent bool, returnCode MqttConnAckCodes) (*MqttConnectAck, error) {
	if header.GetType() != BaseMqtt.CONNACK {
		return nil, fmt.Errorf("Wrong header type [%v] used to initalize connectAck struct", header.GetType())
	}

	return &MqttConnectAck{
		Header:         header,
		SessionPresent: sessionPresent,
		ReturnCode:     returnCode,
	}, nil
}
