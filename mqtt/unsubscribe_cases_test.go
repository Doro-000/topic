package mqtt_test

import (
	Mqtt "github.com/Doro-000/topic/mqtt"
)

var unsubscribe_test_cases PacketTestCase = PacketTestCase{
	"unsubscribe test": {
		packet: &Mqtt.MqttUnsubscribe{
			MqttHeader: Mqtt.MqttHeader{
				PacketType: Mqtt.UNSUBSCRIBE,
			},
			PacketIdentifier: 12,
			Payload:          []string{"filter/1", "filter/+/wildCard", "#"},
		},
		packErr:   nil,
		unpackErr: nil,
	},
}
