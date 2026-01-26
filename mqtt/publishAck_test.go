package mqtt_test

import (
	"bytes"
	"testing"

	Mqtt "github.com/Doro-000/topic/mqtt"
	"github.com/google/go-cmp/cmp"
)

func Test_publishAck_MarshallUnmarshall(t *testing.T) {
	// Arrange
	packet := Mqtt.MqttPublishAck{
		PacketIdentifier: 2,
	}

	// Act

	// encode given packet
	var encodedPacket bytes.Buffer
	marshaller := Mqtt.NewMarshall(&encodedPacket)
	err := packet.Marshall(marshaller)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	// decode given packet
	unmarshaller := Mqtt.NewUnmarshall(&encodedPacket)
	decodedPacket := Mqtt.MqttPublishAck{}
	err = decodedPacket.UnmarshallMqttPublishAck(unmarshaller)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	// Assert
	if !cmp.Equal(decodedPacket, packet) {
		t.Errorf("Packet mismatch after round trip (-want +got):\n%s", cmp.Diff(packet, decodedPacket))
	}
}
