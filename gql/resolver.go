package gql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dikaeinstein/go-graphql-api/data"
	"github.com/dikaeinstein/go-graphql-api/graphqlws"
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
)

// Store describes the data store
type Store interface {
	GetUsersByName(ctx context.Context, name string) ([]data.User, error)
	GetUserByEmail(ctx context.Context, email string) (*data.User, error)
	CreateUser(ctx context.Context, userData data.User) (*data.User, error)
	UpdateUser(context.Context, int, map[string]interface{}) (*data.User, error)
	DeleteUser(ctx context.Context, id int) (*data.User, error)
}

// Resolver resolves the graphql fields
type Resolver struct {
	store  Store
	pubsub graphqlws.PubSub
}

// NewResolver creates a new Resolver
func NewResolver(store Store, pubsub graphqlws.PubSub) *Resolver {
	return &Resolver{store, pubsub}
}

// Users resolves the `users` query
func (r *Resolver) Users(p graphql.ResolveParams) (interface{}, error) {
	ctx, cancelFunc := context.WithTimeout(p.Context, 3*time.Second)
	defer cancelFunc()

	name, ok := p.Args["name"].(string)
	if !ok {
		return nil, nil
	}

	return r.store.GetUsersByName(ctx, name)
}

// User resolves the `user` query
func (r *Resolver) User(p graphql.ResolveParams) (interface{}, error) {
	ctx, cancelFunc := context.WithTimeout(p.Context, 3*time.Second)
	defer cancelFunc()

	email, ok := p.Args["email"].(string)
	if !ok {
		return nil, nil
	}

	user, err := r.store.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// CreateUser resolves the `createUser` mutation
func (r *Resolver) CreateUser(p graphql.ResolveParams) (interface{}, error) {
	ctx, cancelFunc := context.WithTimeout(p.Context, 3*time.Second)
	defer cancelFunc()

	input, ok := p.Args["createUserInput"]
	if !ok {
		return nil, nil
	}

	var u data.User
	mapstructure.Decode(input, &u)
	newUser, err := r.store.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	r.pubsub.Publish("userCreated", newUser)
	return newUser, nil
}

// UpdateUser resolves the `updateUser` mutation
func (r *Resolver) UpdateUser(p graphql.ResolveParams) (interface{}, error) {
	ctx, cancelFunc := context.WithTimeout(p.Context, 3*time.Second)
	defer cancelFunc()

	payload, ok := p.Args["updateUserInput"].(map[string]interface{})
	id, ok2 := p.Args["id"].(int)
	if !ok || !ok2 {
		return nil, nil
	}

	updatedUser, err := r.store.UpdateUser(ctx, id, payload)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return updatedUser, nil
}

// DeleteUser resolves the `deleteUser` mutation
func (r *Resolver) DeleteUser(p graphql.ResolveParams) (interface{}, error) {
	ctx, cancelFunc := context.WithTimeout(p.Context, 3*time.Second)
	defer cancelFunc()

	id, ok := p.Args["id"].(int)
	if !ok {
		return nil, nil
	}

	deletedUser, err := r.store.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return deletedUser, nil
}
