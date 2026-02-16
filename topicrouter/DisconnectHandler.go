package topicrouter

import (
	"fmt"

	mqtt "github.com/Doro-000/topic/mqtt"
	"github.com/Doro-000/topic/topicnetworking"
)

func DisconnectHandler(packet mqtt.GenericPacket, connection topicConnection) error {
	fmt.Printf("Disconnecting client: %s\n", connection.(*topicnetworking.TcpConnection).Client.ClientID)
	err := connection.Close()

	return err
}
