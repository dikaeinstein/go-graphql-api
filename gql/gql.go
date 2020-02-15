package gql

import (
	"github.com/graphql-go/graphql"
)

// Root is the root Query
type Root struct {
	Query *graphql.Object
}

// NewRootQuery creates the Root Query
func NewRootQuery(resolver Resolver) *Root {
	rootQuery := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"users": &graphql.Field{
					Type:        graphql.NewList(graphql.NewNonNull(userType)),
					Description: "Get list of users that match given name",
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type:        graphql.String,
							Description: "Filter users with name",
						},
					},
					Resolve: resolver.Users,
				},
				"user": &graphql.Field{
					Type:        userType,
					Description: "Get user by email",
					Args: graphql.FieldConfigArgument{
						"email": &graphql.ArgumentConfig{
							Type:        graphql.String,
							Description: "Filter by email",
						},
					},
					Resolve: resolver.User,
				},
			},
		},
	)

	return &Root{Query: rootQuery}
}
