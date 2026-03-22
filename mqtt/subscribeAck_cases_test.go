package mqtt_test

import (
	Mqtt "github.com/Doro-000/topic/mqtt"
)

var subscribeAck_test_cases PacketTestCase = PacketTestCase{
	"Subscribe ack test": {
		packet: &Mqtt.MqttSubAck{
			MqttHeader: Mqtt.MqttHeader{
				PacketType: Mqtt.SUBACK,
			},
			PacketIdentifier: 12,
			Payload:          []Mqtt.MqttSubAckCode{Mqtt.SUCCESS_MAX_QOS_0},
		},
		packErr:   nil,
		unpackErr: nil,
	},
}
