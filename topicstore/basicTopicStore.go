package topicstore

type basicTopicStore struct {
	store map[string][]Message
}

func initBasicTopicStore() *basicTopicStore {
	store := make(map[string][]Message, 1000)

	return &basicTopicStore{
		store: store,
	}
}

func (topicStore *basicTopicStore) AddMessages(topicName string, newMessages []Message) error {
	messages, ok := topicStore.store[topicName]

	if ok {
		messages = append(messages, newMessages...)
		topicStore.store[topicName] = messages
	} else {
		topicStore.store[topicName] = newMessages
	}

	return nil
}

func (topicStore *basicTopicStore) GetMessages(topicName string) []Message {
	messages, ok := topicStore.store[topicName]

	if ok {
		return messages
	}
	return []Message{}
}
