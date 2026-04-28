package topicrouter

import (
	"bytes"
	"slices"

	mqtt "github.com/Doro-000/topic/mqtt"
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
)

func UnsubscribeHandler(packet mqtt.GenericPacket, connection topicNetworking.GenericConnection, handlerInput MqttHandlerInput) error {
	unsubPacket := packet.(*mqtt.MqttUnsubscribe)

	clientData := connection.GetClientData()
	session := handlerInput.sessionStore.Get(clientData.TransportId)

	for _, topic := range unsubPacket.Payload {
		session.Subscriptions = slices.DeleteFunc(session.Subscriptions, func(item string) bool {
			return item == topic
		})
	}

	handlerInput.topicStore.RemoveSubscription(clientData.TransportId, unsubPacket.Payload...)

	// respond with unsubAck
	unsubAckHeader, err := mqtt.NewMqttHeader(mqtt.UNSUBACK, false, mqtt.AT_MOST_ONCE, false)
	if err != nil {
		return err
	}

	unsubAckPacket := packetFactory(*unsubAckHeader).(*mqtt.MqttSimplePacket)
	unsubAckPacket.PacketIdentifier = unsubPacket.PacketIdentifier

	var encodedPacket bytes.Buffer
	marshaller := mqtt.NewMarshall(&encodedPacket)
	unsubAckPacket.Marshall(marshaller)

	_, err = connection.Write(encodedPacket.Bytes())

	if err != nil {
		return err
	}
	return nil
}
