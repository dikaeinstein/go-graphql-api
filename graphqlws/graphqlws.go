package graphqlws

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

// ConnectionMessage represents the GraphQL WebSocket message.
type ConnectionMessage struct {
	OperationID string `json:"id,omitempty"`
	Type        string `json:"type"`
	Payload     struct {
		OperationName string                 `json:"operationName,omitempty"`
		Query         string                 `json:"query"`
		Variables     map[string]interface{} `json:"variables"`
	} `json:"payload,omitempty"`
}

type graphqlWS struct {
	eventHandlers       ConnectionEventHandlers
	upgrader            websocket.Upgrader
	subscriptionManager *SubscriptionManager
}

// NewHandler returns a websocket based HTTP handler for graphQL
func NewHandler(u websocket.Upgrader, s *SubscriptionManager, e ConnectionEventHandlers) http.Handler {
	return &graphqlWS{
		eventHandlers:       e,
		upgrader:            u,
		subscriptionManager: s,
	}
}

func (gws *graphqlWS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := gws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to do websocket upgrade: %v", err)
		return
	}
	defer conn.Close()

	connectionACK := map[string]string{
		"type": gqlConnectionAck,
	}
	if err := conn.WriteJSON(connectionACK); err != nil {
		log.Printf("failed to write to ws connection: %v", err)
		return
	}
	handleWSConn(conn, gws.subscriptionManager, gws.eventHandlers)
}

func handleWSConn(conn *websocket.Conn, subMgr *SubscriptionManager, e ConnectionEventHandlers) {
	for {
		var msg ConnectionMessage
		err := conn.ReadJSON(&msg)
		if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			log.Println("connection closed; going away")
			return
		}
		if err != nil {
			log.Println("failed to read websocket message:", err)
			return
		}

		switch msg.Type {
		case gqlStart:
			s := createSubscription(conn, msg, subMgr)
			e.Start(s)
		case gqlStop:
			e.Stop(msg.OperationID)
		case gqlConnectionTerminate:
			e.Close(conn)
		default:
			log.Println("unhandled message", msg.Type)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func createSubscription(conn *websocket.Conn, msg ConnectionMessage, subMgr *SubscriptionManager) *Subscription {
	callback := func(result *graphql.Result) error {
		m := map[string]interface{}{
			"id":      msg.OperationID,
			"type":    gqlData,
			"payload": result,
		}
		if err := conn.WriteJSON(m); err != nil {
			if err == websocket.ErrCloseSent {
				subMgr.RemoveSubscription(msg.OperationID)
				log.Println("subscription removed")
			}
			return errors.Wrap(err, "failed to write to ws connection")
		}

		return nil
	}

	return &Subscription{
		ID:            msg.OperationID,
		RequestString: msg.Payload.Query,
		Variables:     msg.Payload.Variables,
		OperationName: msg.Payload.OperationName,
		Conn:          conn,
		CallBack:      callback,
	}
}
