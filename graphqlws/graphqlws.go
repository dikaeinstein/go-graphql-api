package graphqlws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dikaeinstein/go-graphql-api/pubsub"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ConnectionACKMessage struct {
	OperationID string `json:"id,omitempty"`
	Type        string `json:"type"`
	Payload     struct {
		Query string `json:"query"`
	} `json:"payload,omitempty"`
}

type graphqlWS struct {
	upgrader      websocket.Upgrader
	pubsub        pubsub.PubSub
	eventHandlers ConnectionEventHandlers
}

// NewHandler returns a websocket based HTTP handler for graphQL
func NewHandler(upgrader websocket.Upgrader, pubsub pubsub.PubSub,
	eventHandlers ConnectionEventHandlers) http.Handler {
	return graphqlWS{upgrader, pubsub, eventHandlers}
}

func (gws graphqlWS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := gws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to do websocket upgrade: %v", err)
		return
	}
	connectionACK, err := json.Marshal(map[string]string{
		"type": gqlConnectionAck,
	})
	if err != nil {
		log.Printf("failed to marshal ws connection ack: %v", err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, connectionACK); err != nil {
		log.Printf("failed to write to ws connection: %v", err)
		return
	}
	go handleWSConn(conn, gws.eventHandlers, gws.pubsub)
}

func handleWSConn(conn *websocket.Conn, eventHandlers ConnectionEventHandlers, ps pubsub.PubSub) {
	for {
		_, p, err := conn.ReadMessage()
		if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			log.Println("connection closed")
			return
		}
		if err != nil {
			log.Println("failed to read websocket message:", err)
			return
		}
		var msg ConnectionACKMessage
		if err := json.Unmarshal(p, &msg); err != nil {
			log.Println("failed to unmarshal:", err)
			return
		}
		sub := createSubscriber(conn, msg)
		switch msg.Type {
		case gqlStart:
			eventHandlers.Start(conn, sub, ps)
		case gqlStop:
			eventHandlers.Stop(conn, sub.ID, ps)
		case gqlConnectionTerminate:
			eventHandlers.Close(conn)
		default:
			log.Println("unhandled message", msg.Type)
		}
	}
}

func createSubscriber(conn *websocket.Conn, msg ConnectionACKMessage) *pubsub.Subscriber {
	c := &pubsub.Client{
		Conn:          conn,
		OperationID:   msg.OperationID,
		RequestString: msg.Payload.Query,
	}
	return &pubsub.Subscriber{
		ID:     uuid.New().String(),
		Event:  "userCreated",
		Client: c,
	}
}
