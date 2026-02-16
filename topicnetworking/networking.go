package topicnetworking

import (
	"syscall"
	"time"
)

type ConnectionType int

const (
	RAW_TCP ConnectionType = iota
	WEB_SOCK
)

type TcpListener struct {
	SockFD int
}

type ClientData struct {
	RawAddr        syscall.SockaddrInet4
	ClientID       string // This could reflect the client ID field in the Mqtt protocol
	RemoteAddress  string
	LocalAddress   string
	ConnectionType ConnectionType
	ConnectedAt    time.Time
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

type TcpConnection struct {
	Client   ClientData
	ClientFD int
	PacketBuffer
}
