package topicstore

import mqtt "github.com/Doro-000/topic/mqtt"

type Message struct {
	Qos    mqtt.QoSLevel
	Retain bool
	Data   []byte
}

// TODO: Support wildcards
type TopicStore interface {
	AddMessages(topicName string, newMessages []Message) error
	GetMessages(topicName string) []Message
}

func InitTopicStore() TopicStore {
	basicTopicStore := initBasicTopicStore()

	return basicTopicStore
}
