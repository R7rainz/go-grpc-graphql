package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
	graphserver "github.com/r7rainz/go-grpc-graphql/graphql/internal/graph"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	resolver, err := graphserver.NewResolver(cfg.AccountURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resolver.Close()

	http.Handle("/graphql", handler.NewDefaultServer(resolver.ToExecutableSchema()))
	http.Handle("/playground", playground.Handler("rainz", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
