package main

import (
	"fmt"

	topicEventLoop "github.com/Doro-000/topic/topiceventing"
	topicLog "github.com/Doro-000/topic/topiclog"
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
	topicRouter "github.com/Doro-000/topic/topicrouter"
)

func main() {
	logger := topicLog.NewTopicLogger("MAIN-SERVER", topicLog.LevelFilter{
		Info: true,
		Err:  true,
	})

	// create an event loop
	eventLoop := topicEventLoop.NewEventLoop(false, 10)

	// open a socket for TCP connections
	tcpListener, err := topicNetworking.NewTcpListener()

	if err != nil {
		logger.Error("Failed to create a TCP listener, stopping server!", err)
		return
	}

	newClientHandler := func() error {
		clientFD, clientData, err := tcpListener.Accept()

		if err != nil {
			return err
		}

		connection := topicNetworking.NewTcpConnection(clientFD, clientData)

		err = eventLoop.Add(clientFD, topicEventLoop.Callback{
			Handler: func() error {
				if clientData.ConnectionType == topicNetworking.RAW_TCP {
					err := connection.Fill()
					if err != nil {
						return err
					}
				}

				// Peek data and get packet
				mqttPacket, err := topicRouter.Peek(connection)
				logger.Info(fmt.Sprintf("recieved packet %v", mqttPacket))

				if err != nil {
					return err
				}

				err = topicRouter.RespondTo(mqttPacket, connection)

				if err != nil {
					return err
				}

				return nil
			},
			ErrHandler: func(err error) {
				if err != nil {
					// gracefull shutdown
					connection.Close()
					eventLoop.Remove(connection.ClientFD)
					logger.Info(fmt.Sprintf("Client %s disconnected!", clientData.ClientID))
					logger.Error("Reason: ", err)
				}
			},
		})

		return err
	}

	err = eventLoop.Add(tcpListener.SockFD, topicEventLoop.Callback{
		Handler: newClientHandler,
		ErrHandler: func(err error) {
			logger.Error("Err when accepting connection or registering new client", err)
		},
	})

	if err != nil {
		logger.Error("Failed to register TCP listener with event loop", err)
		return
	}

	err = eventLoop.Wait()

	if err != nil {
		panic(err)
	}
}
