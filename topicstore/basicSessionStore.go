package topicstore

import (
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
)

type basicStore struct {
	store map[string]*Session
}

func initBasicSessionStore() *basicStore {
	hash := make(map[string]*Session)

	return &basicStore{
		store: hash,
	}
}

func (basicStore *basicStore) InitSession(connection topicNetworking.GenericConnection, clientId string, willTopic string, willMessage string, subs []string, keepSession bool) (*Session, error) {
	newSession := Session{}
	clientData := connection.GetClientData()

	newSession.ClientId = clientId
	newSession.WillTopic = willTopic
	newSession.WillMessage = willMessage
	newSession.Subscriptions = subs
	newSession.KeepSession = keepSession
	newSession.Connection = connection

	basicStore.store[clientData.TransportId] = &newSession

	return &newSession, nil
}

func (basicStore *basicStore) Get(transportId string) *Session {
	return basicStore.store[transportId]
}

func (basicStore *basicStore) Delete(transportId string) error {
	delete(basicStore.store, transportId)
	return nil
}

func (basicStore *basicStore) GetAll() []Session {
	res := make([]Session, 0, len(basicStore.store))
	for _, session := range basicStore.store {
		res = append(res, *session)
	}
	return res
}
