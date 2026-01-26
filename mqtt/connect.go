package mqtt

import (
	"errors"
	"fmt"
)

type MqttConnFlags struct {
	Value byte
}

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

type MqttConnPayload struct {
	ClientID    string
	WillTopic   string
	WillMessage string
	Username    string
	Password    string
}

type MqttConnect struct {
	Header    MqttHeader
	Flags     MqttConnFlags
	KeepAlive uint16
	Payload   MqttConnPayload
}

func (packet *MqttConnect) Marshall(marshaller *Marshall) error {
	marshaller.WriteString("MQTT")
	marshaller.WriteByte(0x04)
	marshaller.WriteByte(packet.Flags.Value)
	marshaller.WriteUint16(packet.KeepAlive)
	marshaller.WriteString(packet.Payload.ClientID)

	if packet.Flags.HasWillFlag() {
		marshaller.WriteString(packet.Payload.WillTopic)
		marshaller.WriteString(packet.Payload.WillMessage)
	}

	if packet.Flags.HasUsername() {
		marshaller.WriteString(packet.Payload.Username)
	}

	if packet.Flags.HasPassword() {
		marshaller.WriteString(packet.Payload.Password)
	}

	return marshaller.Error()
}

func (packet *MqttConnect) Unmarshall(unmarshaller *Unmarshall) error {
	protocolName, _ := unmarshaller.ReadString()

	if string(protocolName) != "MQTT" {
		return errors.New("Malformed variable header for connect packet: Wrong protocol name")
	}

	protocolLevel, _ := unmarshaller.ReadByte()

	if protocolLevel != 4 {
		// respond with connack 0x01 (Unacceptable protocol level)
		return fmt.Errorf("Unacceptable protocol level: expected 4 got %v", protocolLevel)
	}

	packet.Flags.Value, _ = unmarshaller.ReadByte()
	packet.KeepAlive = unmarshaller.ReadUint16()
	packet.Payload.ClientID, _ = unmarshaller.ReadString()

	if packet.Flags.HasWillFlag() {
		packet.Payload.WillTopic, _ = unmarshaller.ReadString()
		packet.Payload.WillMessage, _ = unmarshaller.ReadString()
	}

	if packet.Flags.HasUsername() {
		packet.Payload.Username, _ = unmarshaller.ReadString()
	}

	if packet.Flags.HasPassword() {
		packet.Payload.Password, _ = unmarshaller.ReadString()
	}

	return unmarshaller.Error()
}
