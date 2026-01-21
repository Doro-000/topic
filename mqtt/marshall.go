package mqtt

import (
	"fmt"
)

func MarshallMqttPacket(packet GenericPacket) ([]byte, error) {
	payload := packet.GetPayload()
	varHeader := packet.GetVariableHeader()

	varHeaderAndPayloadLength := len(payload) + len(varHeader)
	remainingLength := EncodeRemainingLength(varHeaderAndPayloadLength)

	packetSize := 1 + len(remainingLength) + varHeaderAndPayloadLength
	buf := make([]byte, 0, packetSize)

	buf = append(buf, packet.GetHeader())
	buf = append(buf, remainingLength...)
	buf = append(buf, varHeader...)
	buf = append(buf, payload...)

	if len(buf) != packetSize {
		return []byte{}, fmt.Errorf("Expected packet size: %v, got: %v", packetSize, len(buf))
	}

	return buf, nil
}
