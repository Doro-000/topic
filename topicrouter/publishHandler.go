package topicrouter

import (
	"bytes"

	mqtt "github.com/Doro-000/topic/mqtt"
	topicDataStore "github.com/Doro-000/topic/topicdatastore"
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
)

func publishHandler(packet mqtt.GenericPacket, connection topicNetworking.GenericConnection, handlerInput MqttHandlerInput) error {
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

	responsePacket := mqtt.PacketFactory(*responsePacketHeader).(*mqtt.MqttPublish)
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
		// TODO: submit to worker pool
		go func() {
			conn.Write(encodedPacket.Bytes())
		}()
	}

	if pubPacket.Qos == mqtt.AT_LEAST_ONCE {
		err := sendTo(connection, &mqtt.MqttSimplePacket{
			MqttHeader: mqtt.MqttHeader{
				PacketType: mqtt.PUBACK,
			},
			PacketIdentifier: pubPacket.PacketIdentifier,
		})

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

		pubRecPacket := mqtt.PacketFactory(*pubRecHeader).(*mqtt.MqttSimplePacket)
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
