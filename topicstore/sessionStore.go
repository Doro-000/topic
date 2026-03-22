package topicstore

import (
	topicNetworking "github.com/Doro-000/topic/topicnetworking"
)

// TODO: replace subsciptions with topic store provider ?
// TODO: store QOS 1 and 2 Messages
type Session struct {
	ClientId      string
	WillTopic     string
	WillMessage   string
	Subscriptions []string
	KeepSession   bool
	Connection    topicNetworking.GenericConnection
}

type SessionStore interface {
	InitSession(connection topicNetworking.GenericConnection, clientId string, willTopic string, willMessage string, subs []string, keepSession bool) (*Session, error)
	Get(clientId string) *Session
	Delete(clientId string) error
	GetAll() []Session
}

func InitSessionStore() SessionStore {
	basicStore := initBasicSessionStore()

	return basicStore
}
