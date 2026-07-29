package graph

import (
	gql "github.com/99designs/gqlgen/graphql"
	accountclient "github.com/r7rainz/go-grpc-graphql/account/client"
)

type Resolver struct {
	accountClient *accountclient.Client
}

func NewResolver(accountURL string) (*Resolver, error) {
	accountClient, err := accountclient.NewClient(accountURL)
	if err != nil {
		return nil, err
	}

	return &Resolver{
		accountClient: accountClient,
	}, nil
}

func (r *Resolver) Close() {
	if r.accountClient != nil {
		r.accountClient.Close()
	}
}

func (r *Resolver) ToExecutableSchema() gql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: r,
	})
}
