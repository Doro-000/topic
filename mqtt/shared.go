package mqtt

type MQTTControlPacketType byte

//go:generate stringer -type=MQTTControlPacketType
const (
	RESERVED_PACKET_TYPE MQTTControlPacketType = iota
	CONNECT
	CONNACK
	PUBLISH
	PUBACK
	PUBREC  // Publish received
	PUBREL  // Publish release
	PUBCOMP // Publish complete
	SUBSCRIBE
	SUBACK
	UNSUBSCRIBE
	UNSUBACK
	PINGREQ
	PINGRESP
	DISCONNECT
	RESERVED_PACKET_TYPE_2
)

type QoSLevel byte

const (
	AT_MOST_ONCE QoSLevel = iota
	AT_LEAST_ONCE
	EXACTLY_ONCE
	RESERVED_QOS_LEVEL
)

const MQTT_MAX_PACKET_SIZE = 1500 // 64kB
const MQTT_REMAIN_LEN_MAX = 4     // maximum of 4 bytes for Remaining Length
const MAX_MULTIPLIER_REMAIN_LEN = 128 * 128 * 128

type GenericPacket interface {
	// TODO: add validations to marshalling process
	// TODO: marshall header separately
	Marshall(e *Marshall) error
	Unmarshall(u *Unmarshall) error
}
