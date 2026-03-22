package topicstore

import topicNetworking "github.com/Doro-000/topic/topicnetworking"

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

func (basicStore *basicStore) Get(clientId string) *Session {
	return basicStore.store[clientId]
}

func (basicStore *basicStore) Delete(clientId string) error {
	delete(basicStore.store, clientId)
	return nil
}

func (basicStore *basicStore) GetAll() []Session {
	res := make([]Session, len(basicStore.store))
	for clientId := range basicStore.store {
		res = append(res, *basicStore.store[clientId])
	}
	return res
}
