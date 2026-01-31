package mqtt_test

import (
	"bytes"
	"io"
	"testing"

	Mqtt "github.com/Doro-000/topic/mqtt"
	"github.com/google/go-cmp/cmp"
)

type PacketTestCase = map[string]struct {
	packet    Mqtt.GenericPacket
	packErr   error
	unpackErr error
}
type TestCases = map[Mqtt.MQTTControlPacketType]PacketTestCase

var tests TestCases = TestCases{
	Mqtt.CONNECT:     connect_test_cases,
	Mqtt.CONNACK:     connectAck_test_cases,
	Mqtt.PUBLISH:     publish_test_cases,
	Mqtt.PUBACK:      pubAck_test_cases,
	Mqtt.PUBREC:      pubRec_test_cases,
	Mqtt.PUBREL:      pubRel_test_cases,
	Mqtt.PUBCOMP:     pubComp_test_cases,
	Mqtt.SUBSCRIBE:   nil,
	Mqtt.SUBACK:      nil,
	Mqtt.UNSUBSCRIBE: nil,
	Mqtt.UNSUBACK:    unsubAck_test_cases,
	Mqtt.PINGREQ:     nil,
	Mqtt.PINGRESP:    nil,
	Mqtt.DISCONNECT:  nil,
}

func Test_RoundTrip(t *testing.T) {
	for packetType, testCases := range tests {
		if len(testCases) == 0 {
			continue
		}

		t.Run(packetType.String(), func(t *testing.T) {
			for testName, testCase := range testCases {
				t.Run(testName, func(t *testing.T) {
					// encode given packet
					var encodedPacket bytes.Buffer
					marshaller := Mqtt.NewMarshall(&encodedPacket)
					err := testCase.packet.Marshall(marshaller)
					if err != nil {
						t.Fatalf("%s", err.Error())
					}

					// decode encoded packet
					headerByte, err := encodedPacket.ReadByte()
					if err != nil {
						t.Fatalf("%s", err.Error())
					}

					decodedHeader, err := Mqtt.UnmarshallMqttHeader(headerByte)
					if err != nil {
						t.Fatalf("%s", err.Error())
					}

					decodedPacket := packetFactory(packetType, decodedHeader)

					remainingLen, err := Mqtt.DecodeRemainingLen(&encodedPacket)
					if err != nil {
						t.Fatalf("%s", err.Error())
					}
					unmarshaller := Mqtt.NewUnmarshall(io.LimitReader(&encodedPacket, int64(remainingLen)))

					err = decodedPacket.Unmarshall(unmarshaller)
					if err != nil {
						t.Fatalf("%s", err.Error())
					}

					// Assert
					if !cmp.Equal(decodedPacket, testCase.packet) {
						t.Errorf("Packet mismatch after round trip (-want +got):\n%s", cmp.Diff(testCase.packet, decodedPacket))
					}
				})
			}
		})
	}
}

func packetFactory(packetType Mqtt.MQTTControlPacketType, header *Mqtt.MqttHeader) Mqtt.GenericPacket {
	switch packetType {
	case Mqtt.CONNECT:
		return &Mqtt.MqttConnect{
			MqttHeader: *header,
		}
	case Mqtt.CONNACK:
		return &Mqtt.MqttConnectAck{
			MqttHeader: *header,
		}
	case Mqtt.PUBLISH:
		return &Mqtt.MqttPublish{
			MqttHeader: *header,
		}
	case Mqtt.PUBACK, Mqtt.PUBREC, Mqtt.PUBREL, Mqtt.PUBCOMP, Mqtt.UNSUBACK:
		return &Mqtt.MqttSimplePacket{
			MqttHeader: *header,
		}
	default:
		return nil
	}
}
