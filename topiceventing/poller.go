package topiceventing

import (
	topiclog "github.com/Doro-000/topic/topiclog"
)

type eventLoop struct {
	fd            int
	maxEventsSize int
	logger        *topiclog.TopicLogger
	callbacks     map[int]Callback
}

type Callback struct {
	Handler    func() error
	ErrHandler func(err error)
}

// TODO: poller is not a correct name here
type EventLoop interface {
	Add(fd int, callBack Callback) error
	Remove(fd int) error
	Wait() error
}

const EVENT_TIMEOUT_NSEC int64 = 5 * 1e9
