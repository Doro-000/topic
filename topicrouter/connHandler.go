package topicrouter

import (
	"bytes"

	mqtt "github.com/Doro-000/topic/mqtt"
)

func ConnectHandler(packet mqtt.GenericPacket, connection topicConnection) error {
	connAckHeader, err := mqtt.NewMqttHeader(mqtt.CONNACK, false, mqtt.AT_MOST_ONCE, false)
	if err != nil {
		return err
	}

	connAckPacket := packetFactory(*connAckHeader).(*mqtt.MqttConnectAck)
	connAckPacket.SessionPresent = false
	connAckPacket.ReturnCode = mqtt.ACCEPTED

	var encodedPacket bytes.Buffer
	marshaller := mqtt.NewMarshall(&encodedPacket)
	connAckPacket.Marshall(marshaller)

	_, err = connection.Write(encodedPacket.Bytes())

	if err != nil {
		return err
	}
	return nil
}
