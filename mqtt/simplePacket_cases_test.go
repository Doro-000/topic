package mqtt_test

import (
	Mqtt "github.com/Doro-000/topic/mqtt"
)

var pubAck_test_cases PacketTestCase = PacketTestCase{
	"PUBACK": {
		packet: &Mqtt.MqttSimplePacket{
			MqttHeader: Mqtt.MqttHeader{
				PacketType: Mqtt.PUBACK,
			},
			PacketIdentifier: 10,
		},
		packErr:   nil,
		unpackErr: nil,
	},
}

var pubRec_test_cases PacketTestCase = PacketTestCase{
	"PUBREC": {
		packet: &Mqtt.MqttSimplePacket{
			MqttHeader: Mqtt.MqttHeader{
				PacketType: Mqtt.PUBREC,
			},
			PacketIdentifier: 11,
		},
		packErr:   nil,
		unpackErr: nil,
	},
}

var pubRel_test_cases PacketTestCase = PacketTestCase{
	"PUBREL": {
		packet: &Mqtt.MqttSimplePacket{
			MqttHeader: Mqtt.MqttHeader{
				PacketType: Mqtt.PUBREL,
				Qos:        Mqtt.AT_LEAST_ONCE,
			},
			PacketIdentifier: 12,
		},
		packErr:   nil,
		unpackErr: nil,
	},
}

var pubComp_test_cases PacketTestCase = PacketTestCase{
	"PUBCOMP": {
		packet: &Mqtt.MqttSimplePacket{
			MqttHeader: Mqtt.MqttHeader{
				PacketType: Mqtt.PUBCOMP,
			},
			PacketIdentifier: 13,
		},
		packErr:   nil,
		unpackErr: nil,
	},
}

var unsubAck_test_cases PacketTestCase = PacketTestCase{
	"UNSUBACK": {
		packet: &Mqtt.MqttSimplePacket{
			MqttHeader: Mqtt.MqttHeader{
				PacketType: Mqtt.UNSUBACK,
			},
			PacketIdentifier: 14,
		},
		packErr:   nil,
		unpackErr: nil,
	},
}