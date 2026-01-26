package mqtt_test

import (
	"bytes"
	"testing"

	"github.com/Doro-000/topic/mqtt"
	"github.com/google/go-cmp/cmp"
)

func Test_Publish_MarshallUnmarshall(t *testing.T) {
	// Arrange
	packet := mqtt.MqttPublish{
		TopicName: "a/b",
		Payload:   []byte("hello world"),
	}

	// Act
	var encodedPacket bytes.Buffer
	marshaller := mqtt.NewMarshall(&encodedPacket)
	err := packet.Marshall(marshaller)
	if err != nil {
		t.Fatalf("Failed to marshall packet: %s", err)
	}

	unmarshaller := mqtt.NewUnmarshall(&encodedPacket)
	// The unmarshaller needs the header context to decode correctly
	decodedPacket := mqtt.MqttPublish{Header: packet.Header}
	err = decodedPacket.Unmarshall(unmarshaller)
	if err != nil {
		t.Fatalf("Failed to unmarshall packet: %s", err)
	}

	// Assert
	if !cmp.Equal(decodedPacket, packet) {
		t.Errorf("Packet mismatch after round trip (-want +got):\n%s", cmp.Diff(packet, decodedPacket))
	}
}
