package topicrouter

import (
	"time"

	mqtt "github.com/Doro-000/topic/mqtt"
	"github.com/Doro-000/topic/topicdatastore"
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
)

// TODO: handle will Retain
func connectHandler(packet mqtt.GenericPacket, connection topicNetworking.GenericConnection, handlerInput MqttHandlerInput) error {
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
	err := sendTo(connection, &mqtt.MqttConnectAck{
		SessionPresent: (!cleanSession && currentSession != nil),
		ReturnCode:     mqtt.ACCEPTED,
	})

	return err
}
