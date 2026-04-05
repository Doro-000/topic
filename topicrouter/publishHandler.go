package topicrouter

import (
	"bytes"

	mqtt "github.com/Doro-000/topic/mqtt"
	topicDataStore "github.com/Doro-000/topic/topicdatastore"
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
)

func PublishHandler(packet mqtt.GenericPacket, connection topicNetworking.GenericConnection, handlerInput MqttHandlerInput) error {
	pubPacket := packet.(*mqtt.MqttPublish)
	currentClient := connection.GetClientData()

	// TODO: only save on appropriate qos level
	if pubPacket.Retain {
		message := topicDataStore.Message{
			Qos:  pubPacket.Qos,
			Data: pubPacket.Payload,
		}

		handlerInput.messageStore.AddMessages(pubPacket.TopicName, message)
	}

	var encodedPacket bytes.Buffer
	marshaller := mqtt.NewMarshall(&encodedPacket)
	responsePacketHeader, err := mqtt.NewMqttHeader(mqtt.PUBLISH, false, mqtt.AT_MOST_ONCE, false)
	if err != nil {
		return err
	}

	responsePacket := packetFactory(*responsePacketHeader).(*mqtt.MqttPublish)
	responsePacket.TopicName = pubPacket.TopicName
	responsePacket.Payload = pubPacket.Payload

	responsePacket.Marshall(marshaller)

	clientsToSendTo := make([]topicNetworking.GenericConnection, 0)
	clientsSubscribedToTopic := handlerInput.topicStore.FindClientsSubedToTopic(pubPacket.TopicName)

	for _, client := range clientsSubscribedToTopic {
		clientSession := handlerInput.sessionStore.Get(client)

		if clientSession == nil || clientSession.Connection == nil {
			continue
		}

		clientData := clientSession.Connection.GetClientData()
		if clientData.Connected == false {
			continue
		}

		if clientData.TransportId == currentClient.TransportId {
			continue
		}

		clientsToSendTo = append(clientsToSendTo, clientSession.Connection)
	}

	for _, conn := range clientsToSendTo {
		go func() {
			conn.Write(encodedPacket.Bytes())
		}()
	}

	if pubPacket.Qos == mqtt.AT_LEAST_ONCE {
		// send puback
		pubAckHeader, err := mqtt.NewMqttHeader(mqtt.PUBACK, false, mqtt.AT_MOST_ONCE, false)
		if err != nil {
			return err
		}

		pubAckPacket := packetFactory(*pubAckHeader).(*mqtt.MqttSimplePacket)
		pubAckPacket.PacketIdentifier = pubPacket.PacketIdentifier

		var encodedPacket bytes.Buffer
		marshaller := mqtt.NewMarshall(&encodedPacket)
		pubAckPacket.Marshall(marshaller)

		_, err = connection.Write(encodedPacket.Bytes())

		if err != nil {
			return err
		}
	}

	if pubPacket.Qos == mqtt.EXACTLY_ONCE {
		handlerInput.messageStore.AddUnackPacket(pubPacket.PacketIdentifier)

		// send pubrec
		pubRecHeader, err := mqtt.NewMqttHeader(mqtt.PUBREC, false, mqtt.AT_MOST_ONCE, false)
		if err != nil {
			return err
		}

		pubRecPacket := packetFactory(*pubRecHeader).(*mqtt.MqttSimplePacket)
		pubRecPacket.PacketIdentifier = pubPacket.PacketIdentifier

		var encodedPacket bytes.Buffer
		marshaller := mqtt.NewMarshall(&encodedPacket)
		pubRecPacket.Marshall(marshaller)

		_, err = connection.Write(encodedPacket.Bytes())

		if err != nil {
			return err
		}
	}

	return nil
}

func pubRelHandler(packet mqtt.GenericPacket, connection topicNetworking.GenericConnection, handlerInput MqttHandlerInput) error {
	pubRelPacket := packet.(*mqtt.MqttSimplePacket)

	handlerInput.messageStore.RemoveUnackPacket(pubRelPacket.PacketIdentifier)

	// Send Pubcomp
	pubCompHeader, err := mqtt.NewMqttHeader(mqtt.PUBCOMP, false, mqtt.AT_MOST_ONCE, false)
	if err != nil {
		return err
	}

	pubCompPacket := packetFactory(*pubCompHeader).(*mqtt.MqttSimplePacket)
	pubCompPacket.PacketIdentifier = pubRelPacket.PacketIdentifier

	var encodedPacket bytes.Buffer
	marshaller := mqtt.NewMarshall(&encodedPacket)
	pubCompPacket.Marshall(marshaller)

	_, err = connection.Write(encodedPacket.Bytes())

	if err != nil {
		return err
	}

	return nil

}
