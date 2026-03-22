package topicrouter

import (
	"fmt"

	mqtt "github.com/Doro-000/topic/mqtt"
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
)

func DisconnectHandler(packet mqtt.GenericPacket, connection topicNetworking.GenericConnection, handlerInput MqttHandlerInput) error {
	fmt.Printf("Disconnecting client: %s\n", connection.(*topicNetworking.TcpConnection).Client.TransportId)

	clientData := connection.GetClientData()
	session := handlerInput.sessionStore.Get(clientData.TransportId)

	err := connection.Close()
	session.Connection = nil
	clientData.DisconnectChan <- true
	clientData.KeepAliveTimer.Stop()

	return err
}
