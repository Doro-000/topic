package publish_test

import (
	"bytes"
	"testing"

	BaseMqtt "github.com/Doro-000/topic/mqtt"
	Publish "github.com/Doro-000/topic/mqtt/publish"
	"github.com/google/go-cmp/cmp"
)

func TestMarshallUnmarshall(t *testing.T) {
	// Arrange
	header, err := BaseMqtt.NewMqttHeader(BaseMqtt.PUBLISH, false, 1, false)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	packet := Publish.MqttPublish{
		Header:           *header,
		TopicName:        "test/topic",
		PacketIdentifier: []byte{0, 2},
		Payload:          []byte("testing"),
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

	unmarshalledPacket, err := Publish.UnmarshallMqttPublish(*unmarshalledHeader, incomingPacket)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if !cmp.Equal(*unmarshalledPacket, packet) {
		t.Fail()
	}
}
