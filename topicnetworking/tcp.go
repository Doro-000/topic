package topicnetworking

import (
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type TcpConnection struct {
	Client   ClientData
	ClientFD int
	PacketBuffer
}

func (listener *TcpListener) Accept() (int, ClientData, error) {
	clientFD, client, err := syscall.Accept(listener.SockFD)
	if err != nil {
		return 0, ClientData{}, err
	}

	clientData := ClientData{
		ConnectionType: RAW_TCP,
		ConnectedAt:    time.Now(),
		Connected:      false,
		TimerValue:     DEFAULT_CLIENT_TIMEOUT,
		DisconnectChan: make(chan bool),
	}

	if addr, ok := client.(*syscall.SockaddrInet4); ok {
		clientData.RawAddr = *addr
		clientData.RemoteAddress = fmt.Sprintf("%d.%d.%d.%d", addr.Addr[0], addr.Addr[1], addr.Addr[2], addr.Addr[3])
		clientData.LocalAddress = fmt.Sprintf("%d", addr.Port)
		clientData.TransportId = uuid.New().String()
	}

	fmt.Printf("Accepted Connection: %v\n", clientData)
	return clientFD, clientData, nil
}

func NewTcpListener() (*TcpListener, error) {
	sockFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return nil, err
	}
	err = syscall.SetNonblock(sockFd, true)
	if err != nil {
		return nil, err
	}

	err = syscall.SetsockoptInt(sockFd, syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1)
	if err != nil {
		return nil, err
	}

	err = syscall.Bind(sockFd, &syscall.SockaddrInet4{
		Port: 1883,
		Addr: [4]byte{0, 0, 0, 0},
	})
	if err != nil {
		return nil, err
	}

	err = syscall.Listen(sockFd, syscall.SOMAXCONN)
	if err != nil {
		return nil, err
	}

	return &TcpListener{
		SockFD: sockFd,
	}, nil
}

func (conn *TcpConnection) Fill() error {
	// we have no more space, reset array
	if conn.tail > 0 {
		conn.shift()
	}

	// read from kernel
	n, err := syscall.Read(conn.ClientFD, conn.data[conn.tail:])

	// advance tail
	conn.tail += n

	if err != nil {
		if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
			fmt.Println(fmt.Errorf("Read err, %w", err))
			return nil
		}
		return err
	}

	if n == 0 {
		return errors.New("EOF")
	}

	return nil
}

func (conn *TcpConnection) Read(p []byte) (int, error) {
	if conn.tail == conn.head {
		return 0, nil
	}

	dataAvailable := min(len(p), conn.tail-conn.head)
	copy(p, conn.data[conn.head:conn.head+dataAvailable])
	conn.head += dataAvailable

	return dataAvailable, nil
}

func (conn *TcpConnection) ReadByte() (byte, error) {
	if conn.tail == conn.head {
		return 0, nil
	}

	res := conn.data[0]
	conn.head += 1
	return res, nil
}

func (conn *TcpConnection) Write(p []byte) (int, error) {
	fmt.Printf("Writing to fd: %d | clientData: %v| len: %d\n", conn.ClientFD, conn.Client, len(p))
	n, err := syscall.Write(conn.ClientFD, p)

	if err != nil {
		if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
			fmt.Println(fmt.Errorf("No data yet, %w", err))
			return 0, nil
		}
		return -1, err
	}

	return n, nil
}

func (conn *TcpConnection) Close() error {
	return syscall.Close(conn.ClientFD)
}

func (conn *TcpConnection) GetClientData() *ClientData {
	return &conn.Client
}

func NewTcpConnection(clientFd int, client ClientData) *TcpConnection {
	return &TcpConnection{
		ClientFD: clientFd,
		Client:   client,
		PacketBuffer: PacketBuffer{
			data: make([]byte, 1024),
			head: 0,
			tail: 0,
		},
	}
}
