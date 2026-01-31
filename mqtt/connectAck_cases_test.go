package mqtt_test

import (
	Mqtt "github.com/Doro-000/topic/mqtt"
)

var connectAck_test_cases PacketTestCase = PacketTestCase{
	"Session Present, Accepted": {
		packet: &Mqtt.MqttConnectAck{
			MqttHeader: Mqtt.MqttHeader{
				PacketType: Mqtt.CONNACK,
			},
			SessionPresent: true,
			ReturnCode:     Mqtt.ACCEPTED,
		},
		packErr:   nil,
		unpackErr: nil,
	},
}
