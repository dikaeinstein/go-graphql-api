package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dikaeinstein/go-graphql-api/config"
	"github.com/dikaeinstein/go-graphql-api/data"
	"github.com/dikaeinstein/go-graphql-api/data/postgres"
	"github.com/dikaeinstein/go-graphql-api/gql"
	"github.com/dikaeinstein/go-graphql-api/graphqlws"
	"github.com/dikaeinstein/go-graphql-api/pubsub"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	_ "github.com/lib/pq"
)

// Time to wait before force close on connection.
const closeGracePeriod = 10 * time.Second

func main() {
	cfg := config.New()
	db := connectPostgresDB(cfg)
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	ps := pubsub.NewDefaultPubSub()

	graphql := setupGraphQLHandler(db, ps)
	r.Handle("/graphql", graphql)

	graphqlws := setupGraphQLWSHandler(ps)
	r.Handle("/subscriptions", graphqlws)

	log.Println("Server listening on port:", cfg.Port)
	http.ListenAndServe(":6600", r)
}

func connectPostgresDB(cfg config.Config) *postgres.Postgres {
	connStr := postgres.ConnString(
		cfg.DBName, cfg.DBUser,
		postgres.ConnectTimeout(cfg.DBConnectTimeout),
		postgres.SSLMode("disable"),
	)
	postgresDB, err := postgres.New(connStr)
	if err != nil {
		log.Fatalln(err)
	}

	return postgresDB
}

func setupGraphQLHandler(store data.Store, pubsub pubsub.PubSub) http.Handler {
	resolver := gql.NewResolver(store, pubsub)
	root := gql.NewRoot(resolver)
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:        root.Query,
		Mutation:     root.Mutation,
		Subscription: root.Subscription,
	})
	if err != nil {
		log.Fatal(err)
	}

	return handler.New(&handler.Config{
		Schema:     &schema,
		GraphiQL:   false,
		Playground: true,
		Pretty:     true,
	})
}

func setupGraphQLWSHandler(ps pubsub.PubSub) http.Handler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{"graphql-ws"},
	}
	eventHandlers := graphqlws.ConnectionEventHandlers{
		Start: func(conn *websocket.Conn, sub *pubsub.Subscriber, ps pubsub.PubSub) {
			ps.Subscribe(sub)
		},
		Stop: func(conn *websocket.Conn, subID string, ps pubsub.PubSub) {
			ps.Unsubscribe(subID)
		},
		Close: func(conn *websocket.Conn) {
			time.Sleep(closeGracePeriod)
			conn.Close()
		},
	}
	return graphqlws.NewHandler(upgrader, ps, eventHandlers)
}
