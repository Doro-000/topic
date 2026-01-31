package mqtt_test

import (
	Mqtt "github.com/Doro-000/topic/mqtt"
)

var publish_test_cases PacketTestCase = PacketTestCase{
	"Basic Publish": {
		packet: &Mqtt.MqttPublish{
			MqttHeader: Mqtt.MqttHeader{
				Retain:     true,
				Qos:        Mqtt.AT_LEAST_ONCE,
				Dup:        false,
				PacketType: Mqtt.PUBLISH,
			},
			TopicName: "a/b",
			Payload:   []byte("hello world"),
		},
		packErr:   nil,
		unpackErr: nil,
	},
}
