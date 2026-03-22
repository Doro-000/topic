package topicnetworking

import (
	"io"
	"syscall"
	"time"
)

type ConnectionType int

const (
	RAW_TCP ConnectionType = iota
	WEB_SOCK
)

// If the client doesn't send a CONNECT packet within 30 seconds, we disconnect them
const DEFAULT_CLIENT_TIMEOUT = 30 * time.Second // 30 seconds

type TcpListener struct {
	SockFD int
}

type ClientData struct {
	RawAddr        syscall.SockaddrInet4
	TransportId    string
	RemoteAddress  string
	LocalAddress   string
	ConnectionType ConnectionType

	ConnectedAt        time.Time
	LastPacketRecieved time.Time

	TimerValue     time.Duration
	KeepAliveTimer *time.Timer
	Connected      bool // Flag indicating wheter this client has completed the first MQTT_CONNECT/ACK handshake
	DisconnectChan chan bool
}

type PacketBuffer struct {
	data []byte
	head int
	tail int
}

func (buf *PacketBuffer) shift() {
	unreadData := buf.tail - buf.head
	for i := range unreadData {
		buf.data[i] = buf.data[buf.head+i]
	}
	buf.head, buf.tail = 0, unreadData
}

type GenericConnection interface {
	io.ReadWriteCloser
	io.ByteReader
	GetClientData() *ClientData
}
