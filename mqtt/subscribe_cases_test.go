package mqtt_test

import (
	Mqtt "github.com/Doro-000/topic/mqtt"
)

var subscribe_test_cases PacketTestCase = PacketTestCase{
	"Basic subscribe": {
		packet: &Mqtt.MqttSubscribe{
			MqttHeader: Mqtt.MqttHeader{
				Retain:     false,
				Qos:        Mqtt.AT_LEAST_ONCE,
				Dup:        false,
				PacketType: Mqtt.SUBSCRIBE,
			},
			PacketIdentifier: 12,
			Payload: map[string]Mqtt.QoSLevel{
				"topic1": Mqtt.AT_LEAST_ONCE,
			},
		},
		packErr:   nil,
		unpackErr: nil,
	},
}
