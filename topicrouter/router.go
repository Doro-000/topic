package topicrouter

import (
	"context"
	"fmt"
	"io"

	mqtt "github.com/Doro-000/topic/mqtt"
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
	topicStore "github.com/Doro-000/topic/topicstore"
)

type TopicRouter struct {
	sessionStore    topicStore.SessionStore
	topicStore      topicStore.TopicStore
	handlerRegistry MqttPacketHandlerRegistry
	mainContext     context.Context
}

type MqttHandlerInput struct {
	sessionStore topicStore.SessionStore
	topicStore   topicStore.TopicStore
}

// TODO: return specific error type
type MqttHandlerFunc = func(mqtt.GenericPacket, topicNetworking.GenericConnection, MqttHandlerInput) error
type MqttPacketHandlerRegistry = map[mqtt.MQTTControlPacketType]MqttHandlerFunc

func NewTopicRouter(ctx context.Context, sessionStore topicStore.SessionStore, topicStore topicStore.TopicStore) *TopicRouter {
	return &TopicRouter{
		sessionStore: sessionStore,
		topicStore:   topicStore,
		mainContext:  ctx,
		handlerRegistry: MqttPacketHandlerRegistry{
			mqtt.CONNECT:    ConnectHandler,
			mqtt.PINGREQ:    PingHandler,
			mqtt.DISCONNECT: DisconnectHandler,
			mqtt.SUBSCRIBE:  SubscribeHandler,
			mqtt.PUBLISH:    PublishHandler,
			mqtt.PUBACK: func(mqtt.GenericPacket, topicNetworking.GenericConnection, MqttHandlerInput) error {
				fmt.Print("Hooray!")
				return nil
			},
		},
	}
}

func (router *TopicRouter) RespondTo(packet mqtt.GenericPacket, connection topicNetworking.GenericConnection) error {
	clientData := connection.GetClientData()

	if handler, ok := router.handlerRegistry[packet.GetType()]; ok {
		// First packet sent should be CONNECT
		if clientData.Connected == false && packet.GetType() != mqtt.CONNECT {
			// TODO: return specific error type
			return fmt.Errorf("Client sent %s instead of %s", packet.GetType(), mqtt.CONNECT)
		}

		err := handler(packet, connection, MqttHandlerInput{
			sessionStore: router.sessionStore,
			topicStore:   router.topicStore,
		})

		if err != nil {
			return err
		}

		// do we set the keepAliveTimer to nil ?
		if clientData.KeepAliveTimer != nil {
			clientData.KeepAliveTimer.Reset(clientData.TimerValue)
		}
	} else {
		connection.Close()
		clientData.KeepAliveTimer.Stop()
		return fmt.Errorf("Handler for %s not found!", packet.GetType())
	}
	return nil
}

func Peek(connection topicNetworking.GenericConnection) (mqtt.GenericPacket, error) {
	header, err := connection.ReadByte()
	if err != nil {
		return nil, err
	}
	unmarshalledHeader, err := mqtt.UnmarshallMqttHeader(header)
	if err != nil {
		return nil, err
	}

	remainingLen, err := mqtt.DecodeRemainingLen(connection)
	if err != nil {
		return nil, err
	}

	limitedReader := io.LimitReader(connection, int64(remainingLen))
	unmarshaller := mqtt.NewUnmarshall(limitedReader)

	packet := packetFactory(*unmarshalledHeader)
	err = packet.Unmarshall(unmarshaller)
	if err != nil {
		return nil, err
	}

	return packet, nil
}

func packetFactory(header mqtt.MqttHeader) mqtt.GenericPacket {
	switch header.PacketType {
	case mqtt.CONNECT:
		return &mqtt.MqttConnect{
			MqttHeader: header,
		}
	case mqtt.CONNACK:
		return &mqtt.MqttConnectAck{
			MqttHeader: header,
		}
	case mqtt.PUBLISH:
		return &mqtt.MqttPublish{
			MqttHeader: header,
		}
	case mqtt.PUBACK, mqtt.PUBREC, mqtt.PUBREL, mqtt.PUBCOMP, mqtt.UNSUBACK:
		return &mqtt.MqttSimplePacket{
			MqttHeader: header,
		}
	case mqtt.SUBSCRIBE:
		return &mqtt.MqttSubscribe{
			MqttHeader: header,
		}
	case mqtt.SUBACK:
		return &mqtt.MqttSubAck{
			MqttHeader: header,
		}
	case mqtt.UNSUBSCRIBE:
		return &mqtt.MqttUnsubscribe{
			MqttHeader: header,
		}
	case mqtt.DISCONNECT, mqtt.PINGREQ, mqtt.PINGRESP:
		return &mqtt.MqttHeaderOnlyPacket{
			MqttHeader: header,
		}

	case mqtt.RESERVED_PACKET_TYPE, mqtt.RESERVED_PACKET_TYPE_2:
	default:
		return nil
	}

	return nil
}
