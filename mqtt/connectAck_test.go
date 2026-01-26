package connectack_test

import (
	"bytes"
	"testing"

	BaseMqtt "github.com/Doro-000/topic/mqtt"
	ConnectAck "github.com/Doro-000/topic/mqtt/connectAck"
	"github.com/google/go-cmp/cmp"
)

func TestMarshallUnmarshall(t *testing.T) {
	// Arrange
	header, err := BaseMqtt.NewMqttHeader(BaseMqtt.CONNACK, false, 0, false)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	packet, err := ConnectAck.NewMqttConnectAck(*header, false, ConnectAck.BAD_PROTO_V)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	// Act
	buf, err := BaseMqtt.MarshallMqttPacket(packet)
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

	unmarshalledPacket, err := ConnectAck.UnmarshallMqttConnectAck(*unmarshalledHeader, incomingPacket)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if !cmp.Equal(*unmarshalledPacket, *packet) {
		t.Fail()
	}
}
