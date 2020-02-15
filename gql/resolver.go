package gql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	return r.store.GetUsersByName(ctx, name)
}

// User resolves the `user` query
func (r Resolver) User(p graphql.ResolveParams) (interface{}, error) {
	ctx, cancelFunc := context.WithTimeout(p.Context, 3*time.Second)
	defer cancelFunc()

	email, ok := p.Args["email"].(string)
	if !ok {
		return nil, nil
	}

	user, err := r.store.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}
