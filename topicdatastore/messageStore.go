package topicdatastore

import mqtt "github.com/Doro-000/topic/mqtt"

type Message struct {
	Qos  mqtt.QoSLevel
	Data []byte
}

type MessageStore struct {
	// TODO: [CONFIG] max message per topic
	store map[string][]Message // topic : message[]

	// QOS 2 messages
	unackowledgedPacketIds map[uint16]bool // Packet Id set
}

func NewMessageStore() *MessageStore {
	store := make(map[string][]Message)
	unackowledgedPacketIds := make(map[uint16]bool)

	return &MessageStore{
		store:                  store,
		unackowledgedPacketIds: unackowledgedPacketIds,
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

func (messageStoreInstance *MessageStore) AddUnackPacket(packageId uint16) {
	messageStoreInstance.unackowledgedPacketIds[packageId] = true
}

func (messageStoreInstance *MessageStore) RemoveUnackPacket(packageId uint16) {
	delete(messageStoreInstance.unackowledgedPacketIds, packageId)
}
