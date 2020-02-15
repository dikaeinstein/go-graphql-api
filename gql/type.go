package gql

import (
	"github.com/graphql-go/graphql"
)

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "User",
		Description: "Represents a user",
		Fields: graphql.Fields{
			"id":         &graphql.Field{Type: graphql.Int},
			"name":       &graphql.Field{Type: graphql.String},
			"email":      &graphql.Field{Type: graphql.String},
			"age":        &graphql.Field{Type: graphql.Int},
			"profession": &graphql.Field{Type: graphql.String},
			"friendly":   &graphql.Field{Type: graphql.Boolean},
		},
	},
)
