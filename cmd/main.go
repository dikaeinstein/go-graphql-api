package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dikaeinstein/go-graphql-api/config"
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

	ps := pubsub.NewInMemoryPubSub()

	schema := setupGraphQLSchema(db, ps)
	graphql := setupGraphQLHandler(schema)
	r.Handle("/graphql", graphql)

	graphqlws := setupGraphQLWSHandler(schema, ps)
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

func setupGraphQLSchema(store gql.Store, pubsub graphqlws.PubSub) graphql.Schema {
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

	return schema
}

func setupGraphQLHandler(schema graphql.Schema) http.Handler {
	return handler.New(&handler.Config{
		Schema:     &schema,
		GraphiQL:   false,
		Playground: true,
		Pretty:     true,
	})
}

func setupGraphQLWSHandler(schema graphql.Schema, ps graphqlws.PubSub) http.Handler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{"graphql-ws"},
	}

	subManager := graphqlws.NewSubscriptionManager(&schema, ps)

	eventHandlers := graphqlws.ConnectionEventHandlers{
		Start: func(s *graphqlws.Subscription) {
			subManager.AddSubscription(s)
			log.Println("subscription added")
		},
		Stop: func(subscriptionID string) {
			subManager.RemoveSubscription(subscriptionID)
			log.Println("subscription removed")
		},
		Close: func(conn *websocket.Conn) {
			log.Println("closing graphQL client connection...")
			time.Sleep(closeGracePeriod)
			if err := conn.Close(); err != nil {
				log.Println(err)
			}
			log.Println("connection closed")
		},
	}

	return graphqlws.NewHandler(upgrader, subManager, eventHandlers)
}
