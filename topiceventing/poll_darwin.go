package topiceventing

import (
	"fmt"
	"syscall"

	topiclog "github.com/Doro-000/topic/topiclog"
)

var kevent_timeSpec = syscall.NsecToTimespec(EVENT_TIMEOUT_NSEC)

func newDarwinEventLoop(eventSize int) *EventLoop {
	kernelQueue, err := syscall.Kqueue()
	if err != nil {
		panic(err)
	}

	logger := topiclog.NewTopicLogger("TOPIC-EVENT-LOOP", topiclog.LevelFilter{
		Info: true,
		Err:  true,
	})

	logger.Info("Event loop created")

	return &EventLoop{
		fd:            kernelQueue,
		maxEventsSize: eventSize,
		logger:        logger,
		callbacks:     make(map[int]Callback),
	}
}

func (ev *EventLoop) Add(fd int, callback Callback) error {
	event := syscall.Kevent_t{
		Ident:  uint64(fd),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE,
	}

	ev.callbacks[fd] = callback

	_, err := syscall.Kevent(ev.fd, []syscall.Kevent_t{event}, nil, &kevent_timeSpec)
	return err
}

func (ev *EventLoop) Remove(fd int) error {
	eventToRemove := syscall.Kevent_t{
		Ident: uint64(fd),
		Flags: syscall.EV_DELETE,
	}

	delete(ev.callbacks, fd)

	_, err := syscall.Kevent(ev.fd, []syscall.Kevent_t{eventToRemove}, nil, &kevent_timeSpec)
	return err
}

func (ev *EventLoop) Wait() error {
	events := make([]syscall.Kevent_t, ev.maxEventsSize)

	for {
		numEvents, err := syscall.Kevent(ev.fd, nil, events, &kevent_timeSpec)

		if err != nil {
			if err.Error() == syscall.EINTR.Error() {
				ev.logger.Info("Handled EINTR")
				continue
			}

			ev.logger.Error("Err when waiting for events", err)
			continue
		}

		for idx := range numEvents {
			eventId := events[idx].Ident

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
