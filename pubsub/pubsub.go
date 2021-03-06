package pubsub

import (
	"log"
	"sync"

	"github.com/dikaeinstein/go-graphql-api/graphqlws"
	"github.com/google/uuid"
)

// NewInMemoryPubSub initializes an in-memory pubsub system.
func NewInMemoryPubSub() *InMemoryPubSub {
	return &InMemoryPubSub{subscribers: sync.Map{}}
}

// InMemoryPubSub implements the PubSub interface with an in-memory map.
type InMemoryPubSub struct {
	subscribers sync.Map
}

// Subscribe registers the given handler for the event.
func (ps *InMemoryPubSub) Subscribe(event string, handler graphqlws.Handler) string {
	s := newSubscriber(event, handler)
	ps.subscribers.Store(s.ID, s)
	return s.ID
}

// Publish publishes to all subscribers of the given event.
func (ps *InMemoryPubSub) Publish(event string, payload interface{}) {
	ps.subscribers.Range(func(k, v interface{}) bool {
		subscriber, ok := v.(Subscriber)
		if !ok {
			return true
		}
		if subscriber.Event == event {
			if err := subscriber.Handler(payload); err != nil {
				log.Println(err)
				return true
			}
		}
		return true
	})
}

// Unsubscribe removes the subscriber with given subID.
func (ps *InMemoryPubSub) Unsubscribe(subID string) {
	ps.subscribers.Delete(subID)
}

// Subscriber represents a client that wants to subscribe to an event.
// You must specify the handler that will be called when the event fires.
type Subscriber struct {
	ID      string
	Event   string
	Handler graphqlws.Handler
}

// newSubscriber creates a new instance of a Subscriber.
func newSubscriber(event string, handler graphqlws.Handler) Subscriber {
	return Subscriber{
		ID:      uuid.New().String(),
		Event:   event,
		Handler: handler,
	}
}
