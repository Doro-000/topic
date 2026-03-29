package topicdatastore

import (
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
)

type Session struct {
	ClientId      string
	WillTopic     string
	WillMessage   string
	Subscriptions []string
	KeepSession   bool
	Connection    topicNetworking.GenericConnection
}

type SessionStore struct {
	store              map[string]*Session // TransportID: session
	persistentSessions map[string]*Session // ClientID: session
}

func NewSessionStore() *SessionStore {
	sessionStore := make(map[string]*Session)
	persistentSessionStore := make(map[string]*Session)

	return &SessionStore{
		store:              sessionStore,
		persistentSessions: persistentSessionStore,
	}
}

func (storeInstance *SessionStore) InitSession(newSession Session) error {
	clientData := newSession.Connection.GetClientData()

	storeInstance.store[clientData.TransportId] = &newSession

	return nil
}

func (storeInstance *SessionStore) Get(transportId string) *Session {
	return storeInstance.store[transportId]
}

func (storeInstance *SessionStore) Delete(transportId string) error {
	delete(storeInstance.store, transportId)
	return nil
}

func (storeInstance *SessionStore) GetAll() []Session {
	res := make([]Session, 0, len(storeInstance.store))
	for _, session := range storeInstance.store {
		res = append(res, *session)
	}
	return res
}

func (storeInstance *SessionStore) GetPersistedSessionForClient(clientId string) *Session {
	oldSession, ok := storeInstance.persistentSessions[clientId]
	if ok {
		return oldSession
	}
	return nil
}

func (storeInstance *SessionStore) RemoveOldSessionForClient(clientId string) {
	delete(storeInstance.persistentSessions, clientId)
}

func (storeInstance *SessionStore) RestoreSession(connection topicNetworking.GenericConnection, persistedSession *Session) error {
	// Remove from persisted sessions
	delete(storeInstance.persistentSessions, persistedSession.ClientId)

	// add to live sessions
	persistedSession.Connection = connection

	clientData := connection.GetClientData()
	storeInstance.store[clientData.TransportId] = persistedSession
	return nil
}

func (storeInstance *SessionStore) PersistSession(disconnectedSession *Session) {
	// Remove from live sessions
	delete(storeInstance.store, disconnectedSession.ClientId)

	// move to peristed sessions
	disconnectedSession.Connection = nil
	storeInstance.persistentSessions[disconnectedSession.ClientId] = disconnectedSession
}
