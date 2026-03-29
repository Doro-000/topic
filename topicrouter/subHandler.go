package topicrouter

import (
	"bytes"

	mqtt "github.com/Doro-000/topic/mqtt"
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
)

func SubscribeHandler(packet mqtt.GenericPacket, connection topicNetworking.GenericConnection, handlerInput MqttHandlerInput) error {
	subPacket := packet.(*mqtt.MqttSubscribe)

	clientData := connection.GetClientData()
	session := handlerInput.sessionStore.Get(clientData.TransportId)

	// topic: qosLevel
	for topic := range subPacket.Payload {
		session.Subscriptions = append(session.Subscriptions, topic)
		handlerInput.topicStore.AddSubscription(topic, clientData.TransportId)
	}

	// respond with subAck
	subAckHeader, err := mqtt.NewMqttHeader(mqtt.SUBACK, false, mqtt.AT_MOST_ONCE, false)
	if err != nil {
		return err
	}

	subAckPacket := packetFactory(*subAckHeader).(*mqtt.MqttSubAck)
	subAckPacket.PacketIdentifier = subPacket.PacketIdentifier
	subAckPacket.Payload = make([]mqtt.MqttSubAckCode, len(subPacket.Payload))

	for idx := range len(subAckPacket.Payload) {
		subAckPacket.Payload[idx] = mqtt.SUCCESS_MAX_QOS_0
	}

	var encodedPacket bytes.Buffer
	marshaller := mqtt.NewMarshall(&encodedPacket)
	subAckPacket.Marshall(marshaller)

	_, err = connection.Write(encodedPacket.Bytes())

	if err != nil {
		return err
	}
	return nil
}
