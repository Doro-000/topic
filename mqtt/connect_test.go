package mqtt_test

import (
	"bytes"
	"testing"

	Mqtt "github.com/Doro-000/topic/mqtt"
	"github.com/google/go-cmp/cmp"
)

func Test_connect_MarshallUnmarshall(t *testing.T) {
	// Arrange
	packet := Mqtt.MqttConnect{
		Flags: Mqtt.MqttConnFlags{
			Value: 194,
		},
		KeepAlive: 1,
		Payload: Mqtt.MqttConnPayload{
			ClientID: "something",
			Username: "Me",
			Password: "Pass",
		},
	}

	// Act

	// encode given packet
	var encodedPacket bytes.Buffer
	marshaller := Mqtt.NewMarshall(&encodedPacket)
	err := packet.Marshall(marshaller)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	// decode encoded packet
	unmarshaller := Mqtt.NewUnmarshall(&encodedPacket)
	decodedPacket := Mqtt.MqttConnect{}
	err = decodedPacket.Unmarshall(unmarshaller)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	// Assert
	if !cmp.Equal(decodedPacket, packet) {
		t.Fail()
	}
}
