package topicrouter

import (
	"bytes"
	"time"

	mqtt "github.com/Doro-000/topic/mqtt"
	"github.com/Doro-000/topic/topicdatastore"
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
)

// TODO: handle will Retain
func ConnectHandler(packet mqtt.GenericPacket, connection topicNetworking.GenericConnection, handlerInput MqttHandlerInput) error {
	connPacket := packet.(*mqtt.MqttConnect)
	sessionStore := handlerInput.sessionStore
	clientData := connection.GetClientData()

	// Retrive old session
	oldSession := sessionStore.GetPersistedSessionForClient(connPacket.ClientID)

	currentSession := handlerInput.sessionStore.Get(clientData.TransportId)
	cleanSession := connPacket.MqttConnFlags.CleanSession

	if cleanSession {
		if oldSession != nil {
			sessionStore.RemoveOldSessionForClient(connPacket.ClientID)
		}

		sessionStore.InitSession(topicdatastore.Session{
			ClientId: connPacket.ClientID, WillTopic: connPacket.WillTopic, WillMessage: connPacket.WillMessage, Subscriptions: []string{}, KeepSession: false, Connection: connection,
		})
	} else {
		if oldSession != nil {
			if connPacket.HasWillFlag {
				oldSession.WillMessage = connPacket.WillMessage
				oldSession.WillTopic = connPacket.WillTopic
			} else {
				oldSession.WillMessage = ""
				oldSession.WillTopic = ""
			}
			sessionStore.RestoreSession(connection, oldSession)
		} else {
			sessionStore.InitSession(topicdatastore.Session{
				ClientId: connPacket.ClientID, WillTopic: connPacket.WillTopic, WillMessage: connPacket.WillMessage, Subscriptions: []string{}, KeepSession: true, Connection: connection,
			})
		}
	}

	clientData.Connected = true
	clientData.LastPacketRecieved = time.Now()

	if connPacket.KeepAlive != 0 {
		// time.Duration needs nanoseconds, keepalive is in seconds
		keepAliveNano := float64(connPacket.KeepAlive) * 1e+9 * mqtt.MQTT_KEEP_ALIVE_TIMEOUT_FACTOR
		clientData.TimerValue = time.Duration(keepAliveNano)
	}

	clientData.KeepAliveTimer = time.NewTimer(clientData.TimerValue)

	// monitor timer
	go func() {
		select {
		case <-clientData.KeepAliveTimer.C:
			// Client timed out
			connection.Close()
			// clean up?
		case <-clientData.DisconnectChan:
			// Client gracefully disonnected
			clientData.KeepAliveTimer.Stop()
			clientData.KeepAliveTimer = nil
			return
		}
	}()

	// RESPOND With ACK
	connAckHeader, err := mqtt.NewMqttHeader(mqtt.CONNACK, false, mqtt.AT_MOST_ONCE, false)
	if err != nil {
		return err
	}

	connAckPacket := packetFactory(*connAckHeader).(*mqtt.MqttConnectAck)
	connAckPacket.SessionPresent = (!cleanSession && currentSession != nil)
	connAckPacket.ReturnCode = mqtt.ACCEPTED

	var encodedPacket bytes.Buffer
	marshaller := mqtt.NewMarshall(&encodedPacket)
	connAckPacket.Marshall(marshaller)

	_, err = connection.Write(encodedPacket.Bytes())

	if err != nil {
		return err
	}
	return nil
}
