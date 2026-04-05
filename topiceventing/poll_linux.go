package topiceventing

import (
	"fmt"
	"syscall"

	topiclog "github.com/Doro-000/topic/topiclog"
)

var epoll_timeout = int(EVENT_TIMEOUT_NSEC / 1e+6)

func NewEventLoop(eventSize int) *eventLoop {
	epollInstance, err := syscall.EpollCreate1(0)
	if err != nil {
		panic(err)
	}

	logger := topiclog.NewTopicLogger("TOPIC-EVENT-LOOP", topiclog.LevelFilter{
		Info: true,
		Err:  true,
	})

	logger.Info("Event loop created")

	return &eventLoop{
		fd:            epollInstance,
		maxEventsSize: eventSize,
		logger:        logger,
		callbacks:     make(map[int]Callback),
	}
}

func (ev *eventLoop) Add(fd int, callback Callback) error {
	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN | syscall.EPOLLET,
		Fd:     int32(fd),
	}

	ev.callbacks[fd] = callback

	err := syscall.EpollCtl(ev.fd, syscall.EPOLL_CTL_ADD, fd, &event)
	return err
}

func (ev *eventLoop) Remove(fd int) error {
	eventToRemove := syscall.EpollEvent{
		Events: syscall.EPOLLOUT,
		Fd:     int32(fd),
	}

	delete(ev.callbacks, fd)

	err := syscall.EpollCtl(ev.fd, syscall.EPOLL_CTL_DEL, fd, &eventToRemove)
	return err
}

func (ev *eventLoop) Wait() error {
	events := make([]syscall.EpollEvent, ev.maxEventsSize)

	for {
		numEvents, err := syscall.EpollWait(ev.fd, events, epoll_timeout)

		if err != nil {
			if err.Error() == syscall.EINTR.Error() {
				ev.logger.Info("Handled EINTR")
				continue
			}

			ev.logger.Error("Err when waiting for events", err)
			continue
		}

		for idx := range numEvents {
			eventId := events[idx].Fd

			if callback, ok := ev.callbacks[int(eventId)]; ok {
				err := callback.Handler()
				if err != nil {
					callback.ErrHandler(err)
				}
			} else {
				ev.logger.Error(fmt.Sprintf("Call back for event %v not found", events[idx]), nil)
			}

		}
	}
}
