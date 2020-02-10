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
					Type: graphql.NewList(
						graphql.NewNonNull(User.Type()),
					),
					Description: "List of users that match given name",
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type:        graphql.String,
							Description: "Filter users with name",
						},
					},
					Resolve: resolver.Users,
				},
			},
		},
	)

	return &Root{Query: rootQuery}
}
