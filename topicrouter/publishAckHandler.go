package topicrouter

import (
	mqtt "github.com/Doro-000/topic/mqtt"
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
)

func pubAckHandler(packet mqtt.GenericPacket, connection topicNetworking.GenericConnection, handlerInput MqttHandlerInput) error {
	return nil
}
