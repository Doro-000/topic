package main

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"
	"time"

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

	mainContext, mainCancel := context.WithCancel(context.Background())

	// create router
	router := topicRouter.NewTopicRouter(mainContext)

	// open a socket for TCP connections
	tcpListener, err := topicNetworking.NewTcpListener()

	if err != nil {
		logger.Error("Failed to create a TCP listener, stopping server!", err)
		return
	}

	// create an event loop
	eventLoop := topicEventLoop.NewEventLoop(10)

	// Client Handler
	newClientHandler := func() error {
		clientFD, clientData, err := tcpListener.Accept()

		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("Accepted Connection: %s | proto: %d", clientData.TransportId, clientData.ConnectionType))
		connection := topicNetworking.NewTcpConnection(clientFD, clientData)

		// Disconnect client if they don't send anything for 30sec
		deadlineContext, signalClientConnected := context.WithTimeout(mainContext, 30*time.Second)

		go inactiveClientDisconnect(*connection, eventLoop, *logger, deadlineContext, signalClientConnected)

		// Register handler for this client
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

				if err != nil {
					return err
				}

				logger.Info(fmt.Sprintf("recieved packet %v", mqttPacket.GetType()))

				signalClientConnected()
				err = router.RespondTo(mqttPacket, connection)

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

					// publish will message
					err := router.ClientFuneral(*connection)

					if err != nil {
						logger.Error("Failed to publish will message: ", err)
					}

					logger.Info(fmt.Sprintf("Client %s disconnected!", clientData.TransportId))
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

	ctx, stop := signal.NotifyContext(mainContext, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		err = eventLoop.Wait()
		if err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()

	fmt.Println("Shutting down server...")
	mainCancel()

	syscall.Close(tcpListener.SockFD)

}

func inactiveClientDisconnect(connection topicNetworking.TcpConnection, eventLoop topicEventLoop.EventLoop, logger topicLog.TopicLogger, deadlineContext context.Context, cancelFunc context.CancelFunc) {
	<-deadlineContext.Done()
	switch deadlineContext.Err() {
	case context.DeadlineExceeded:
		// close connection
		connection.Close()

		// remove handler
		eventLoop.Remove(connection.ClientFD)

		logger.Info(fmt.Sprintf("Client %s disconnected!", connection.Client.TransportId))
		logger.Error("Reason: ", errors.New("Timed out before sending packet"))
	case context.Canceled:
		logger.Info("here")
		cancelFunc()
	}
}
