package gql

import (
	"github.com/graphql-go/graphql"
)

const (
	// User graphql schema
	User = user("user")
)

type user string

func (user) Type() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "User",
			Description: "Represents a user",
			Fields: graphql.Fields{
				"id":         &graphql.Field{Type: graphql.Int},
				"name":       &graphql.Field{Type: graphql.String},
				"age":        &graphql.Field{Type: graphql.Int},
				"profession": &graphql.Field{Type: graphql.String},
				"friendly":   &graphql.Field{Type: graphql.Boolean},
			},
		},
	)
}
