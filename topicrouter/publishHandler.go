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
	responsePacketHeader, err := mqtt.NewMqttHeader(mqtt.PUBLISH, false, mqtt.AT_LEAST_ONCE, false)
	if err != nil {
		return err
	}

	responsePacket := packetFactory(*responsePacketHeader).(*mqtt.MqttPublish)
	responsePacket.TopicName = pubPacket.TopicName
	responsePacket.Payload = pubPacket.Payload

	responsePacket.Marshall(marshaller)

	clientsSubscribedToTopic := handlerInput.topicStore.FindClientsSubedToTopic(pubPacket.TopicName)

	for _, client := range clientsSubscribedToTopic {
		clientSession := handlerInput.sessionStore.Get(client)

		if clientSession == nil {
			continue
		}

		if clientSession.Connection == nil {
			continue
		}

		clientData := clientSession.Connection.GetClientData()
		if clientData.Connected == false {
			continue
		}

		if clientData.TransportId == currentClient.TransportId {
			continue
		}

		go func() {
			clientSession.Connection.Write(encodedPacket.Bytes())
		}()

	}

	return nil
}
