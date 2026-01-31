package mqtt_test

import (
	Mqtt "github.com/Doro-000/topic/mqtt"
)

var connect_test_cases PacketTestCase = PacketTestCase{
	"Basic connect packet": {
		packet: &Mqtt.MqttConnect{
			MqttHeader: Mqtt.MqttHeader{
				PacketType: Mqtt.CONNECT,
			},
			ProtocolName:  "MQTT",
			ProtocolLevel: 4,
			MqttConnFlags: Mqtt.MqttConnFlags{
				CleanSession:  false,
				HasWillFlag:   false,
				HasWillQoS:    Mqtt.AT_LEAST_ONCE,
				HasWillRetain: false,
				HasPassword:   true,
				HasUsername:   true,
			},
			KeepAlive: 1,
			MqttConnPayload: Mqtt.MqttConnPayload{
				ClientID: "something",
				Username: "Me",
				Password: "Pass",
			},
		},
		packErr:   nil,
		unpackErr: nil,
	},
}
