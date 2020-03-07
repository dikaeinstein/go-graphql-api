package gql

import (
	"github.com/graphql-go/graphql"
)

// Root contains the root Query, Mutation and Subscription.
type Root struct {
	Query        *graphql.Object
	Mutation     *graphql.Object
	Subscription *graphql.Object
}

// NewRoot initializes the root query, mutation and subscription.
func NewRoot(resolver *Resolver) *Root {
	return &Root{
		Query:        newRootQuery(resolver),
		Mutation:     newRootMutation(resolver),
		Subscription: newRootSubscription(resolver),
	}
}

func newRootQuery(resolver *Resolver) *graphql.Object {
	return graphql.NewObject(
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
}

func newRootMutation(resolver *Resolver) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"createUser": &graphql.Field{
					Name:        "createUser",
					Description: "Creates a new user and returns the user ID",
					Type:        userType, // the return type for this field
					Args: graphql.FieldConfigArgument{
						"createUserInput": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(createUserInput),
						},
					},
					Resolve: resolver.CreateUser,
				},
				"updateUser": &graphql.Field{
					Name:        "updateUser",
					Description: "Updates user that matches `id` with given payload",
					Type:        userType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.Int),
						},
						"updateUserInput": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(updateUserInput),
						},
					},
					Resolve: resolver.UpdateUser,
				},
				"deleteUser": &graphql.Field{
					Name:        "deleteUser",
					Description: "Deletes user from the data store",
					Type:        userType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.Int),
						},
					},
					Resolve: resolver.DeleteUser,
				},
			},
		},
	)
}

func newRootSubscription(resolver *Resolver) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Subscription",
			Fields: graphql.Fields{
				"userCreated": &graphql.Field{
					Name:        "userCreated",
					Description: "Subscribe to userCreated events",
					Type:        userType,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return p.Info.RootValue, nil
					},
				},
			},
		},
	)
}
