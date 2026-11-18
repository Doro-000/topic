package topicrouter

import (
	"bytes"

	mqtt "github.com/Doro-000/topic/mqtt"
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
)

func pubRelHandler(packet mqtt.GenericPacket, connection topicNetworking.GenericConnection, handlerInput MqttHandlerInput) error {
	pubRelPacket := packet.(*mqtt.MqttSimplePacket)

	handlerInput.messageStore.RemoveUnackPacket(pubRelPacket.PacketIdentifier)

	// Send Pubcomp
	pubCompHeader, err := mqtt.NewMqttHeader(mqtt.PUBCOMP, false, mqtt.AT_MOST_ONCE, false)
	if err != nil {
		return err
	}

	pubCompPacket := mqtt.PacketFactory(*pubCompHeader).(*mqtt.MqttSimplePacket)
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
