package data

import "context"

// User shape
type User struct {
	ID         int
	Name       string
	Age        int
	Profession string
	Friendly   bool
}

// Store describes the data store
type Store interface {
	GetUsersByName(ctx context.Context, name string) ([]User, error)
}
