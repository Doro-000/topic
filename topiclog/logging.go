package topiclog

import (
	"fmt"
	"time"
)

type LevelFilter struct {
	Info  bool
	Debug bool
	Warn  bool
	Err   bool
	Fatal bool
}

type TopicLogger struct {
	// dest        int // TODO: introduce custom logging destinations
	source      string
	levelfilter LevelFilter
}

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string, err error)
	Fatal(msg string)
}

func NewTopicLogger(src string, lvlFilter LevelFilter) *TopicLogger {
	return &TopicLogger{
		source:      src,
		levelfilter: lvlFilter,
	}
}

func (l *TopicLogger) Info(msg string) {
	if l.levelfilter.Info {
		now := time.Now()
		// TODO: standardize logging colors
		out := fmt.Sprintf("\033[94mINFO\033[0m | %v | %s | %s", now, l.source, msg)

		fmt.Println(out)
	}

}

func (l *TopicLogger) Error(msg string, err error) {
	if l.levelfilter.Err {
		now := time.Now()
		// TODO: standardize logging colors
		out := fmt.Sprintf("\033[91mERROR\033[0m | %v | %s | %s", now, l.source, msg)

		fmt.Println(out)
		fmt.Printf("\t\t %s\n", err.Error())
	}
}
