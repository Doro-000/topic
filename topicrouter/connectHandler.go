package topicrouter

import (
	"bytes"
	"time"

	mqtt "github.com/Doro-000/topic/mqtt"
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
)

// TODO: handle will Retain
func ConnectHandler(packet mqtt.GenericPacket, connection topicNetworking.GenericConnection, handlerInput MqttHandlerInput) error {
	connPacket := packet.(*mqtt.MqttConnect)
	sessionStore := handlerInput.sessionStore

	currentSession := handlerInput.sessionStore.Get(connPacket.ClientID)
	cleanSession := connPacket.MqttConnFlags.CleanSession

	if cleanSession {
		if currentSession != nil {
			sessionStore.Delete(connPacket.ClientID)
		}

		sessionStore.InitSession(connection, connPacket.ClientID, connPacket.WillTopic, connPacket.WillMessage, []string{}, false)
	} else {
		if currentSession != nil {
			if connPacket.HasWillFlag {
				currentSession.WillMessage = connPacket.WillMessage
				currentSession.WillTopic = connPacket.WillTopic
			} else {
				currentSession.WillMessage = ""
				currentSession.WillTopic = ""
			}
		} else {
			sessionStore.InitSession(connection, connPacket.ClientID, connPacket.WillTopic, connPacket.WillMessage, []string{}, true)
		}
	}

	clientData := connection.GetClientData()
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
