package connect

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	BaseMqtt "github.com/Doro-000/topic/mqtt"
)

type MqttConnFlags struct {
	Value byte
}

type MqttConnPayload struct {
	ClientID    string
	WillTopic   string
	WillMessage string
	Username    string
	Password    string
}

type MqttConnect struct {
	Header    BaseMqtt.MqttHeader
	Flags     MqttConnFlags
	KeepAlive []byte
	Payload   MqttConnPayload
}

// connect flag getters
func (f *MqttConnFlags) GetCleanSession() bool {
	return (f.Value>>1)&0x01 == 1
}

func (f *MqttConnFlags) HasWillFlag() bool {
	return (f.Value>>2)&0x01 == 1
}

func (f *MqttConnFlags) HasWillQoS() byte {
	return (f.Value >> 3) & 0x03
}

func (f *MqttConnFlags) HasWillRetain() bool {
	return (f.Value>>5)&0x01 == 1
}

func (f *MqttConnFlags) HasPassword() bool {
	return (f.Value>>6)&0x01 == 1
}

func (f *MqttConnFlags) HasUsername() bool {
	return (f.Value>>7)&0x01 == 1
}

func (f *MqttConnect) GetHeader() byte {
	return f.Header.Value
}

func (f *MqttConnect) GetVariableHeader() []byte {
	// length of proto name
	varHeader := []byte{0x00, 0x04}

	protoName := []byte("MQTT")
	varHeader = append(varHeader, protoName...)

	// protocol level
	varHeader = append(varHeader, 0x04)
	varHeader = append(varHeader, f.Flags.Value)
	varHeader = append(varHeader, f.KeepAlive...)

	return varHeader
}

func (f *MqttConnect) GetPayload() []byte {
	payload := []byte{}
	// clientID len
	payload = binary.BigEndian.AppendUint16(payload, uint16(len(f.Payload.ClientID)))
	payload = append(payload, []byte(f.Payload.ClientID)...)

	if f.Flags.HasWillFlag() {
		payload = binary.BigEndian.AppendUint16(payload, uint16(len(f.Payload.WillTopic)))
		payload = append(payload, []byte(f.Payload.WillTopic)...)

		payload = binary.BigEndian.AppendUint16(payload, uint16(len(f.Payload.WillMessage)))
		payload = append(payload, []byte(f.Payload.WillMessage)...)
	}

	if f.Flags.HasUsername() {
		payload = binary.BigEndian.AppendUint16(payload, uint16(len(f.Payload.Username)))
		payload = append(payload, []byte(f.Payload.Username)...)
	}

	if f.Flags.HasPassword() {
		payload = binary.BigEndian.AppendUint16(payload, uint16(len(f.Payload.Password)))
		payload = append(payload, []byte(f.Payload.Password)...)
	}
	return payload
}

func UnmarshallMqttConnect(header BaseMqtt.MqttHeader, packet io.Reader) (*MqttConnect, error) {
	if header.GetType() != BaseMqtt.CONNECT {
		return nil, fmt.Errorf("Called Connect unmarshall on packet type: %v", header.GetType())
	}

	res := MqttConnect{Header: header}

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

	protocolName := packetUnmarshall.String()
	protocolLevel := packetUnmarshall.Uint8()

	if string(protocolName) != "MQTT" {
		return nil, errors.New("Malformed variable header for connect packet")
	}

	if protocolLevel != 4 {
		// respond with connack 0x01 (Unacceptable protocol level)
		return nil, errors.New("Unacceptable protocol level")
	}

	res.Flags.Value = packetUnmarshall.Uint8()

	res.KeepAlive = []byte{packetUnmarshall.Uint8(), packetUnmarshall.Uint8()}

	res.Payload.ClientID = packetUnmarshall.String()

	if res.Flags.HasWillFlag() {
		res.Payload.WillTopic = packetUnmarshall.String()
		res.Payload.WillMessage = packetUnmarshall.String()
	}

	if res.Flags.HasUsername() {
		res.Payload.Username = packetUnmarshall.String()
	}

	if res.Flags.HasPassword() {
		res.Payload.Password = packetUnmarshall.String()
	}

	if packetUnmarshall.Error() != nil {
		return nil, packetUnmarshall.Error()
	}

	return &res, nil
}
