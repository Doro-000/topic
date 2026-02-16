package topiceventing

import (
	"runtime"

	topiclog "github.com/Doro-000/topic/topiclog"
)

type EventLoop struct {
	fd            int
	maxEventsSize int
	logger        *topiclog.TopicLogger
	callbacks     map[int]Callback
}

type Callback struct {
	Handler    func() error
	ErrHandler func(err error)
}

type Poller interface {
	Add(fd int, callBack Callback) error
	Remove(fd int) error
	Wait() error
}

const EVENT_TIMEOUT_NSEC int64 = 5 * 1e9

func NewEventLoop(useDefault bool, eventSize int) Poller {
	if useDefault {
		// TODO: expose a default go lang event loop(somehow)
		return nil
	}

	// TODO: work on IOCP if there is time
	// TODO: implement event loop using io_uring
	if runtime.GOOS == "windows" || runtime.GOOS == "linux" {
		return nil
	}

	return newDarwinEventLoop(eventSize)
}
