package topicdatastore

import mqtt "github.com/Doro-000/topic/mqtt"

type Message struct {
	Qos  mqtt.QoSLevel
	Data []byte
}

type MessageStore struct {
	// TODO: [CONFIG] max message per topic
	store map[string][]Message // topic : message[]
}

func NewMessageStore() *MessageStore {
	store := make(map[string][]Message)

	return &MessageStore{
		store: store,
	}
}

func (messageStoreInstance *MessageStore) AddMessages(topicName string, newMessages ...Message) error {
	messages, ok := messageStoreInstance.store[topicName]

	if ok {
		messages = append(messages, newMessages...)
		messageStoreInstance.store[topicName] = messages
	} else {
		messageStoreInstance.store[topicName] = newMessages
	}

	return nil
}

func (messageStoreInstance *MessageStore) GetMessages(topicName string) []Message {
	messages, ok := messageStoreInstance.store[topicName]

	if ok {
		return messages
	}
	return []Message{}
}
