package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/r7rainz/go-grpc-graphql/account/internal/repository"
	"github.com/r7rainz/go-grpc-graphql/account/internal/service"
)

type AppConfig struct {
	DatabaseURL string `envconfig:"ACCOUNT_DATABASE_URL" required:"true"`
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
	_ = accountService

	log.Println("account service connected to database")
}
