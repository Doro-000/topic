package mqtt

import (
	"errors"
	"fmt"
)

type MqttConnFlags struct {
	CleanSession  bool
	HasWillFlag   bool
	HasWillQoS    QoSLevel
	HasWillRetain bool
	HasPassword   bool
	HasUsername   bool
}

func NewConnectFlag(val byte) (*MqttConnFlags, error) {
	res := MqttConnFlags{}

	if (val>>1)&0x01 == 1 {
		res.CleanSession = true
	}

	if (val>>2)&0x01 == 1 {
		res.HasWillFlag = true
	}

	if qos := QoSLevel((val >> 3) & 0x03); qos != RESERVED_QOS_LEVEL {
		res.HasWillQoS = qos
	} else {
		return nil, fmt.Errorf("Reserved QOS level used!")
	}

	if (val>>5)&0x01 == 1 {
		res.HasWillRetain = true
	}

	if (val>>6)&0x01 == 1 {
		res.HasPassword = true
	}

	if (val>>7)&0x01 == 1 {
		res.HasUsername = true
	}

	return &res, nil
}

func (f *MqttConnFlags) getFlagByte() (val byte) {
	if f.CleanSession {
		val |= (1 << 1)
	}

	if f.HasWillFlag {
		val |= (1 << 2)
	}

	if f.HasWillQoS != 0 {
		val |= (byte(f.HasWillQoS) << 3)
	}

	if f.HasWillRetain {
		val |= (1 << 5)
	}

	if f.HasPassword {
		val |= (1 << 6)
	}

	if f.HasUsername {
		val |= (1 << 7)
	}

	return
}

type MqttConnPayload struct {
	ClientID    string
	WillTopic   string
	WillMessage string
	Username    string
	Password    string
}

type MqttConnect struct {
	MqttHeader
	ProtocolName  string
	ProtocolLevel byte
	MqttConnFlags
	KeepAlive uint16
	MqttConnPayload
}

func (p *MqttConnect) getPayloadLen() (val int) {
	val += (2 + len(p.ClientID))

	if p.HasWillFlag {
		val += (2 + len(p.WillTopic))
		val += (2 + len(p.WillMessage))
	}

	if p.HasUsername {
		val += (2 + len(p.Username))
	}

	if p.HasPassword {
		val += (2 + len(p.Password))
	}

	return
}

func (packet *MqttConnect) Marshall(marshaller *Marshall) error {
	/** Header */
	packet.MqttHeader.Marshall(marshaller)

	/** Remaining Length */
	// 10 => length of the variable header
	remainingLength := EncodeRemainingLength(10 + packet.getPayloadLen())
	marshaller.WriteBytes(remainingLength)

	/** Variable Header */

	// protocol name
	marshaller.WriteString("MQTT")

	// protocol level
	marshaller.WriteByte(0x04)

	// connect flags
	marshaller.WriteByte(packet.getFlagByte())

	// keep alive
	marshaller.WriteUint16(packet.KeepAlive)
	marshaller.WriteString(packet.ClientID)

	/** Payload */
	if packet.HasWillFlag {
		marshaller.WriteString(packet.WillTopic)
		marshaller.WriteString(packet.WillMessage)
	}

	if packet.HasUsername {
		marshaller.WriteString(packet.Username)
	}

	if packet.HasPassword {
		marshaller.WriteString(packet.Password)
	}

	return marshaller.Error()
}

func (packet *MqttConnect) Unmarshall(unmarshaller *Unmarshall) error {
	packet.ProtocolName = unmarshaller.ReadString()

	if string(packet.ProtocolName) != "MQTT" {
		return errors.New("Malformed variable header for connect packet: Wrong protocol name")
	}

	packet.ProtocolLevel = unmarshaller.ReadByte()

	if packet.ProtocolLevel != 4 {
		// respond with connack 0x01 (Unacceptable protocol level)
		return fmt.Errorf("Unacceptable protocol level: expected 4 got %v", packet.ProtocolLevel)
	}

	connectFlag := unmarshaller.ReadByte()
	flag, err := NewConnectFlag(connectFlag)

	if err != nil {
		return err
	}
	packet.MqttConnFlags = *flag

	packet.KeepAlive = unmarshaller.ReadUint16()
	packet.ClientID = unmarshaller.ReadString()

	if packet.HasWillFlag {
		packet.WillTopic = unmarshaller.ReadString()
		packet.WillMessage = unmarshaller.ReadString()
	}

	if packet.HasUsername {
		packet.Username = unmarshaller.ReadString()
	}

	if packet.HasPassword {
		packet.Password = unmarshaller.ReadString()
	}

	return unmarshaller.Error()
}
