package topicrouter

import (
	"bytes"

	mqtt "github.com/Doro-000/topic/mqtt"
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
)

func pingHandler(packet mqtt.GenericPacket, connection topicNetworking.GenericConnection, input MqttHandlerInput) error {
	pingRespHeader, err := mqtt.NewMqttHeader(mqtt.PINGRESP, false, mqtt.AT_MOST_ONCE, false)
	if err != nil {
		return err
	}

	pingRespPacket := mqtt.PacketFactory(*pingRespHeader).(*mqtt.MqttHeaderOnlyPacket)

	var encodedPacket bytes.Buffer
	marshaller := mqtt.NewMarshall(&encodedPacket)
	pingRespPacket.Marshall(marshaller)

	_, err = connection.Write(encodedPacket.Bytes())

	if err != nil {
		return err
	}

	input.loggerInstance.Info("Responded with PINGRESP")

	return nil
}
