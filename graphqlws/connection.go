package graphqlws

import (
	"github.com/dikaeinstein/go-graphql-api/pubsub"
	"github.com/gorilla/websocket"
)

const (
	// Constants for operation message types
	gqlConnectionInit      = "connection_init"
	gqlConnectionAck       = "connection_ack"
	gqlConnectionKeepAlive = "ka"
	gqlConnectionError     = "connection_error"
	gqlConnectionTerminate = "connection_terminate"
	gqlStart               = "start"
	gqlData                = "data"
	gqlError               = "error"
	gqlComplete            = "complete"
	gqlStop                = "stop"
)

type ConnectionEventHandlers struct {
	// Close is called whenever the connection is closed, regardless of
	// whether this happens because of an error or a deliberate termination
	// by the client.
	Close func(conn *websocket.Conn)

	// Start handler is called whenever the client demands that a GraphQL
	// operation be started (typically a subscription). Event handlers
	// are expected to take the necessary steps to register the operation
	// and send data back to the client with the results eventually.
	Start func(conn *websocket.Conn, sub *pubsub.Subscriber, pubsub pubsub.PubSub)

	// Stop handler is called whenever the client stops a previously
	// started GraphQL operation (typically a subscription). Event handlers
	// are expected to unregister the operation and stop sending result
	// data to the client.
	Stop func(conn *websocket.Conn, subID string, pubsub pubsub.PubSub)
}
