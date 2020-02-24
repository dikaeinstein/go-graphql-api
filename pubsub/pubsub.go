package pubsub

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// PubSub is the interface that describes the publish and subscribe system
type PubSub interface {
	Subscribe(s *Subscriber)
	Publish(event string, payload interface{})
	Unsubscribe(subID string)
}

// NewDefaultPubSub initializes an in-memory pubsub system
func NewDefaultPubSub() PubSub {
	return &DefaultPubSub{subscribers: sync.Map{}}
}

// Subscriber wraps the GRAPHQL client and event that it subscribes to
type Subscriber struct {
	ID     string
	Event  string
	Client *Client
}

// Client represents the GRAPHQL client that wants to subscribe to an event.
type Client struct {
	Conn          *websocket.Conn
	RequestString string
	OperationID   string
}

type DefaultPubSub struct {
	subscribers sync.Map
}

func (ps *DefaultPubSub) Subscribe(s *Subscriber) {
	ps.subscribers.Store(s.ID, s)
}

func (ps *DefaultPubSub) Publish(event string, payload interface{}) {
	ps.subscribers.Range(func(k, v interface{}) bool {
		subscriber, ok := v.(*Subscriber)
		if !ok {
			return true
		}
		message, err := json.Marshal(map[string]interface{}{
			"type":    "data",
			"id":      subscriber.Client.OperationID,
			"payload": payload,
		})
		if err != nil {
			log.Printf("failed to marshal message: %v", err)
			return true
		}
		if err := subscriber.Client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			if err == websocket.ErrCloseSent {
				ps.Unsubscribe(k.(string))
				return true
			}
			log.Printf("failed to write to ws connection: %v", err)
			return true
		}
		return true
	})
}

func (ps *DefaultPubSub) Unsubscribe(subID string) {
	ps.subscribers.Delete(subID)
}
