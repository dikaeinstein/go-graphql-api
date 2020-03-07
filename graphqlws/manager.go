package graphqlws

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/pkg/errors"
)

// Handler represents the handler func that should be triggered when an event fires.
type Handler func(payload interface{}) error

// PubSub is the interface that describes the publish and subscribe system.
type PubSub interface {
	Subscribe(event string, handler Handler) (subID string)
	Publish(event string, payload interface{})
	Unsubscribe(subID string)
}

// SubscriptionManager manages the graphQL subscriptions.
type SubscriptionManager struct {
	PubSub        PubSub
	Schema        *graphql.Schema
	subscriptions map[string]*Subscription
}

// CallBack is executed when an event is fired.
type CallBack func(*graphql.Result) error

// Subscription represents the graphQL client subscription.
type Subscription struct {
	ID            string
	RequestString string
	Variables     map[string]interface{}
	OperationName string
	Conn          *websocket.Conn
	CallBack      CallBack
	SubscriberID  string
}

// NewSubscriptionManager creates a new subscription manager.
func NewSubscriptionManager(schema *graphql.Schema, ps PubSub) *SubscriptionManager {
	return &SubscriptionManager{
		Schema:        schema,
		PubSub:        ps,
		subscriptions: make(map[string]*Subscription),
	}
}

// AddSubscription adds a new subscription to the subscription manager.
func (sm *SubscriptionManager) AddSubscription(s *Subscription) error {
	source := source.NewSource(&source.Source{
		Body: []byte(s.RequestString),
		Name: "GraphQL subscription request",
	})
	document, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		return errors.Wrap(err, "failed to parse subscription query")
	}
	validation := graphql.ValidateDocument(sm.Schema, document, graphql.SpecifiedRules)
	if !validation.IsValid {
		return fmt.Errorf("subscription query validation failed: %#v", validation.Errors)
	}

	var subscriptionName string
	var args map[string]interface{}
	for _, node := range document.Definitions {
		if node.GetKind() == "OperationDefinition" {
			def, _ := node.(*ast.OperationDefinition)
			rootField, _ := def.GetSelectionSet().Selections[0].(*ast.Field)
			subscriptionName = rootField.Name.Value

			fields := sm.Schema.SubscriptionType().Fields()
			args, err = getArgumentValues(fields[subscriptionName].Args, rootField.Arguments, s.Variables)
			break
		}
	}

	sID := sm.PubSub.Subscribe("userCreated", func(payload interface{}) error {
		result := graphql.Execute(graphql.ExecuteParams{
			Schema:        *sm.Schema,
			AST:           document,
			OperationName: s.OperationName,
			Args:          args,
			Root:          payload,
		})
		return s.CallBack(result)
	})

	// add new subscription
	s.SubscriberID = sID
	sm.subscriptions[s.ID] = s
	return nil
}

// RemoveSubscription removes the a previously added subscription.
func (sm *SubscriptionManager) RemoveSubscription(subscriptionID string) {
	subID := sm.subscriptions[subscriptionID].SubscriberID
	sm.PubSub.Unsubscribe(subID)
	// delete subscription
	delete(sm.subscriptions, subscriptionID)
}
