package main

import (
	"log"
	"net/http"

	"github.com/dikaeinstein/go-graphql-api/config"
	"github.com/dikaeinstein/go-graphql-api/data"
	"github.com/dikaeinstein/go-graphql-api/data/postgres"
	"github.com/dikaeinstein/go-graphql-api/gql"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.New()
	db := connectPostgresDB(cfg)
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	graphql := setupGraphQLHandler(db)
	r.Handle("/graphql", graphql)

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

func setupGraphQLHandler(store data.Store) http.Handler {
	resolver := gql.NewResolver(store)
	root := gql.NewRootQuery(resolver)
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: root.Query})
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
