package connect_test

import (
	"bytes"
	"testing"

	BaseMqtt "github.com/Doro-000/topic/mqtt"
	Connect "github.com/Doro-000/topic/mqtt/connect"
	"github.com/google/go-cmp/cmp"
)

func TestMarshallUnmarshall(t *testing.T) {
	// Arrange
	header, err := BaseMqtt.NewMqttHeader(BaseMqtt.CONNECT, false, 0, false)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	packet := Connect.MqttConnect{
		Header: *header,
		Flags: Connect.MqttConnFlags{
			Value: 194,
		},
		KeepAlive: []byte{0, 1},
		Payload: Connect.MqttConnPayload{
			ClientID: "something",
			Username: "Me",
			Password: "Pass",
		},
	}

	// Act
	buf, err := BaseMqtt.MarshallMqttPacket(&packet)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	// Assert
	incomingPacket := bytes.NewBuffer(buf)
	headerByte, err := incomingPacket.ReadByte()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	unmarshalledHeader, err := BaseMqtt.UnmarshallMqttHeader(headerByte)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	unmarshalledPacket, err := Connect.UnmarshallMqttConnect(*unmarshalledHeader, incomingPacket)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if !cmp.Equal(*unmarshalledPacket, packet) {
		t.Fail()
	}
}
