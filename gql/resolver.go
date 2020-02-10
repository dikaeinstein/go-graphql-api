package gql

import (
	"context"
	"time"

	"github.com/dikaeinstein/go-graphql-api/data"
	"github.com/graphql-go/graphql"
)

// Resolver resolves the graphql fields
type Resolver struct {
	store data.Store
}

// NewResolver creates a new Resolver
func NewResolver(store data.Store) Resolver {
	return Resolver{store}
}

// Users resolves the `users` query
func (r Resolver) Users(p graphql.ResolveParams) (interface{}, error) {
	ctx, cancelFunc := context.WithTimeout(p.Context, 3*time.Second)
	defer cancelFunc()

	name, ok := p.Args["name"].(string)
	if !ok {
		return nil, nil
	}

	users, err := r.store.GetUsersByName(ctx, name)
	if err != nil {
		return users, err
	}

	return users, nil
}
