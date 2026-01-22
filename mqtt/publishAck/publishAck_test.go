package publishack_test

import (
	"bytes"
	"testing"

	BaseMqtt "github.com/Doro-000/topic/mqtt"
	PublishAck "github.com/Doro-000/topic/mqtt/publishAck"
	"github.com/google/go-cmp/cmp"
)

func TestMarshallUnmarshall(t *testing.T) {
	// Arrange
	header, err := BaseMqtt.NewMqttHeader(BaseMqtt.PUBACK, false, 0, false)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	packet := PublishAck.MqttPublishAck{
		Header:           *header,
		PacketIdentifier: []byte{0, 2},
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

	unmarshalledPacket, err := PublishAck.UnmarshallMqttPublishAck(*unmarshalledHeader, incomingPacket)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if !cmp.Equal(*unmarshalledPacket, packet) {
		t.Fail()
	}
}
