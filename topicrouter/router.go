package topicrouter

import (
	"fmt"
	"io"

	mqtt "github.com/Doro-000/topic/mqtt"
)

type topicConnection interface {
	io.ReadWriteCloser
	io.ByteReader
}

type MqttHandlerFunc = func(mqtt.GenericPacket, topicConnection) error
type MqttPacketHandlerRegistry = map[mqtt.MQTTControlPacketType]MqttHandlerFunc

var Registry = MqttPacketHandlerRegistry{
	mqtt.CONNECT:    ConnectHandler,
	mqtt.PINGREQ:    PingHandler,
	mqtt.DISCONNECT: DisconnectHandler,
}

func RespondTo(packet mqtt.GenericPacket, connection topicConnection) error {
	if handler, ok := Registry[packet.GetType()]; ok {
		err := handler(packet, connection)

		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Handler for %s not found!", packet.GetType())
	}
	return nil
}

func Peek(connection topicConnection) (mqtt.GenericPacket, error) {
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
