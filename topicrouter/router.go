package topicrouter

import (
	"bytes"
	"context"
	"fmt"
	"io"

	mqtt "github.com/Doro-000/topic/mqtt"
	topicDataStore "github.com/Doro-000/topic/topicdatastore"
	topicLog "github.com/Doro-000/topic/topiclog"
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
)

type TopicRouter struct {
	topicStore      *topicDataStore.TopicStore
	sessionStore    *topicDataStore.SessionStore
	messageStore    *topicDataStore.MessageStore
	handlerRegistry MqttPacketHandlerRegistry
	mainContext     context.Context
	loggerInstance  *topicLog.TopicLogger
}

type MqttHandlerInput struct {
	topicStore     *topicDataStore.TopicStore
	sessionStore   *topicDataStore.SessionStore
	messageStore   *topicDataStore.MessageStore
	loggerInstance *topicLog.TopicLogger
}

// TODO: return specific error type
type MqttHandlerFunc = func(mqtt.GenericPacket, topicNetworking.GenericConnection, MqttHandlerInput) error
type MqttPacketHandlerRegistry = map[mqtt.MQTTControlPacketType]MqttHandlerFunc

func NewTopicRouter(ctx context.Context) *TopicRouter {
	// Storage
	sessionStore := topicDataStore.NewSessionStore()
	messageStore := topicDataStore.NewMessageStore()
	topicStore := topicDataStore.NewTopicStore()

	logger := topicLog.NewTopicLogger("ROUTER", topicLog.LevelFilter{
		Info: true,
		Err:  true,
	})

	return &TopicRouter{
		mainContext: ctx,

		sessionStore: sessionStore,
		topicStore:   topicStore,
		messageStore: messageStore,

		loggerInstance: logger,

		handlerRegistry: MqttPacketHandlerRegistry{
			mqtt.CONNECT:     connectHandler,
			mqtt.PINGREQ:     pingHandler,
			mqtt.DISCONNECT:  disconnectHandler,
			mqtt.SUBSCRIBE:   subscribeHandler,
			mqtt.PUBLISH:     publishHandler,
			mqtt.PUBACK:      pubAckHandler,
			mqtt.PUBREL:      pubRelHandler,
			mqtt.UNSUBSCRIBE: unsubscribeHandler,
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
			sessionStore:   router.sessionStore,
			messageStore:   router.messageStore,
			topicStore:     router.topicStore,
			loggerInstance: router.loggerInstance,
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

func (router *TopicRouter) ClientFuneral(closedConnection topicNetworking.TcpConnection) error {
	session := router.sessionStore.Get(closedConnection.Client.TransportId)

	if session.WillTopic == "" || session.WillMessage == "" {
		return nil
	}

	// publish will message
	var encodedPacket bytes.Buffer
	marshaller := mqtt.NewMarshall(&encodedPacket)
	responsePacketHeader, err := mqtt.NewMqttHeader(mqtt.PUBLISH, false, mqtt.AT_MOST_ONCE, false)
	if err != nil {
		return err
	}

	responsePacket := mqtt.PacketFactory(*responsePacketHeader).(*mqtt.MqttPublish)
	responsePacket.TopicName = session.WillTopic
	responsePacket.Payload = []byte(session.WillMessage)

	responsePacket.Marshall(marshaller)

	clientsSubscribedToTopic := router.topicStore.FindClientsSubedToTopic(session.WillTopic)

	clientsToSendTo := make([]topicNetworking.GenericConnection, 0)
	for _, client := range clientsSubscribedToTopic {
		clientSession := router.sessionStore.Get(client)

		if clientSession == nil || clientSession.Connection == nil {
			continue
		}

		clientData := clientSession.Connection.GetClientData()
		if clientData.Connected == false {
			continue
		}

		if clientData.TransportId == closedConnection.Client.TransportId {
			continue
		}

		clientsToSendTo = append(clientsToSendTo, clientSession.Connection)
	}

	for _, conn := range clientsToSendTo {
		go func() {
			conn.Write(encodedPacket.Bytes())
		}()
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

	packet, err := mqtt.PacketFactory(unmarshalledHeader.PacketType)
	err = packet.Unmarshall(unmarshaller)
	if err != nil {
		return nil, err
	}

	return packet, nil
}

func sendTo(conn topicNetworking.GenericConnection, packet mqtt.GenericPacket) error {
	packet, err := mqtt.PacketFactory(packet.GetType())
	if err != nil {
		return err
	}

	var encodedPacket bytes.Buffer
	marshaller := mqtt.NewMarshall(&encodedPacket)
	packet.Marshall(marshaller)

	_, err = conn.Write(encodedPacket.Bytes())

	return err
}
