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

var createUserInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name:        "CreateUserInput",
		Description: "CreateUserInput represents arguments passed to createUser mutation",
		Fields: graphql.InputObjectConfigFieldMap{
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String)},
			"email": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String)},
			"age": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int)},
			"profession": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String)},
			"friendly": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Boolean)},
		},
	},
)

var updateUserInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name:        "UpdateUserInput",
		Description: "UpdateUserInput represents arguments passed to updateUser mutation",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":       &graphql.InputObjectFieldConfig{Type: graphql.String},
			"email":      &graphql.InputObjectFieldConfig{Type: graphql.String},
			"age":        &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"profession": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"friendly":   &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
		},
	},
)
