package topicdatastore

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

type set = map[string]bool
type node struct {
	children map[string]*node
	word     string
	clients  set
}

type TopicStore struct {
	trieRoot *node
	sysTrie  *node
	// Maybe have more info like trie size
}

func newNode(word string) *node {
	return &node{
		children: make(map[string]*node),
		clients:  make(set),
		word:     word,
	}
}

func (store *TopicStore) AddSubscription(topic string, clientId string) error {
	err := validateTopic(topic)
	if err != nil {
		return err
	}

	words := strings.Split(topic, "/")

	currNode := store.trieRoot
	for _, word := range words {
		if existingNode, ok := currNode.children[word]; ok {
			currNode = existingNode
		} else {
			newNode := newNode(word)
			currNode.children[word] = newNode
			currNode = newNode
		}
	}

	currNode.clients[clientId] = true
	return nil
}

func (store *TopicStore) RemoveSubscription(clientId string, topics ...string) {
	for _, topic := range topics {
		words := strings.Split(topic, "/")

		currNode := store.trieRoot
		for _, word := range words {
			if currWordNode, ok := currNode.children[word]; ok && currNode.word == word {
				currNode = currWordNode
			} else {
				// topic no longer exists, break and go to next topic
				// perhaps we should we accumulate errors for these
				break
			}
		}
		delete(currNode.children, clientId)
	}
}

// TODO: rewrite with iterative method
func (store *TopicStore) FindClientsSubedToTopic(topic string) []string {
	words := strings.Split(topic, "/")

	clients := make([]string, 0)
	if len(words) == 0 {
		return clients
	}

	walkTrie(store.trieRoot, words, 0, &clients)
	return clients
}

func walkTrie(node *node, words []string, depth int, clients *[]string) {
	if child, ok := node.children["#"]; ok {
		*clients = append(*clients, slices.Collect(maps.Keys(child.clients))...)
	}

	// recursion guard
	// we've reached the last word
	if depth == len(words) {
		*clients = append(*clients, slices.Collect(maps.Keys(node.clients))...)
		return
	}

	// single level match
	if child, ok := node.children["+"]; ok {
		walkTrie(child, words, depth+1, clients)
	}

	// exact match
	if child, ok := node.children[words[depth]]; ok {
		walkTrie(child, words, depth+1, clients)
	}
}

func validateTopic(topic string) error {
	if len(topic) == 0 {
		return fmt.Errorf("Malformed Topic: %s", topic)
	}

	// Dont support topics starting with '.' as that's our root and could cause edge cases
	if topic[0] == '.' {
		// TODO: have proper error types
		return errors.New("Topic can't start with a .")
	}

	topicLen := len(topic)
	// Check for Max length (65535)

	hierarchyCount := 0
	for i := range topicLen {
		currChar := topic[i]

		if currChar == '#' {
			isOnlyChar := topicLen == 1
			hasSlashBehind := i != 0 && topic[i-1] == '/'
			hasNothingAfter := i == topicLen-1

			if !(isOnlyChar || (hasSlashBehind && hasNothingAfter)) {
				return fmt.Errorf("Malformed Topic: %s", topic)
			}
		}

		if currChar == '+' {
			beforeOk := i == 0 || topic[i-1] == '/'
			afterOk := i == topicLen-1 || topic[i+1] == '/'

			if !(beforeOk && afterOk) {
				return fmt.Errorf("Malformed Topic: %s", topic)
			}
		}

		if currChar == '/' {
			hierarchyCount++
		}
	}

	// Check for Max hierarchLimit

	return nil
}

func NewTopicStore() *TopicStore {
	root := newNode(".")
	sysTrie := newNode("$SYS")

	return &TopicStore{
		trieRoot: root,
		sysTrie:  sysTrie,
	}
}
