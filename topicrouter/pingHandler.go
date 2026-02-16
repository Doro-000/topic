package topicrouter

import (
	"bytes"
	"fmt"

	mqtt "github.com/Doro-000/topic/mqtt"
)

func PingHandler(packet mqtt.GenericPacket, connection topicConnection) error {
	fmt.Print("Responding to ping\n")
	pingRespHeader, err := mqtt.NewMqttHeader(mqtt.PINGRESP, false, mqtt.AT_MOST_ONCE, false)
	if err != nil {
		return err
	}

	pingRespPacket := packetFactory(*pingRespHeader).(*mqtt.MqttHeaderOnlyPacket)

	var encodedPacket bytes.Buffer
	marshaller := mqtt.NewMarshall(&encodedPacket)
	pingRespPacket.Marshall(marshaller)

	_, err = connection.Write(encodedPacket.Bytes())

	if err != nil {
		return err
	}
	return nil
}
