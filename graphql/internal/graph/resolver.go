package graph

import (
	gql "github.com/99designs/gqlgen/graphql"
	"github.com/r7rainz/go-grpc-graphql/graphql/internal/clients/account"
	"github.com/r7rainz/go-grpc-graphql/graphql/internal/clients/catalog"
	"github.com/r7rainz/go-grpc-graphql/graphql/internal/clients/order"
)

type Resolver struct {
	accountClient *account.Client
	catalogClient *catalog.Client
	orderClient   *order.Client
}

func NewResolver(accountURL, catalogURL, orderURL string) (*Resolver, error) {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return nil, err
	}

	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return nil, err
	}

	orderClient, err := order.NewClient(orderURL)
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return nil, err
	}

	return &Resolver{
		accountClient: accountClient,
		catalogClient: catalogClient,
		orderClient:   orderClient,
	}, nil
}

func (r *Resolver) Close() {
	if r.accountClient != nil {
		r.accountClient.Close()
	}
	if r.catalogClient != nil {
		r.catalogClient.Close()
	}
	if r.orderClient != nil {
		r.orderClient.Close()
	}
}

func (r *Resolver) ToExecutableSchema() gql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: r,
	})
}
