package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/r7rainz/go-grpc-graphql/account/internal/repository"
	"github.com/r7rainz/go-grpc-graphql/account/internal/server"
	"github.com/r7rainz/go-grpc-graphql/account/internal/service"
)

type AppConfig struct {
	DatabaseURL string `envconfig:"ACCOUNT_DATABASE_URL" required:"true"`
	GRPCPort    int    `envconfig:"ACCOUNT_GRPC_PORT" default:"8081"`
}

func main() {
	var cfg AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	repo, err := repository.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	accountService := service.NewService(repo)

	log.Printf("account gRPC server listening on :%d", cfg.GRPCPort)
	if err := server.ListenGRPC(accountService, cfg.GRPCPort); err != nil {
		log.Fatal(err)
	}
}
