package topicrouter

import (
	"bytes"
	"slices"

	mqtt "github.com/Doro-000/topic/mqtt"
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
	topicStore "github.com/Doro-000/topic/topicstore"
)

func PublishHandler(packet mqtt.GenericPacket, connection topicNetworking.GenericConnection, handlerInput MqttHandlerInput) error {
	pubPacket := packet.(*mqtt.MqttPublish)
	currentClient := connection.GetClientData()

	// TODO: only save on appropriate qos level
	if pubPacket.Retain {
		message := topicStore.Message{
			Qos:    pubPacket.Qos,
			Retain: pubPacket.Retain,
			Data:   pubPacket.Payload,
		}

		handlerInput.topicStore.AddMessages(pubPacket.TopicName, []topicStore.Message{message})
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

	// find the subscription for each active session
	// TODO: we should have Get By topic, or use the topic store to resolve all subscriptionss ?
	for _, client := range handlerInput.sessionStore.GetAll() {
		if client.Connection == nil {
			continue
		}
		clientData := client.Connection.GetClientData()

		if clientData.TransportId == currentClient.TransportId {
			continue
		}

		if clientData.Connected == false {
			continue
		}

		if slices.Contains(client.Subscriptions, pubPacket.TopicName) {
			go func() {
				client.Connection.Write(encodedPacket.Bytes())
			}()
		}
	}

	return nil
}
