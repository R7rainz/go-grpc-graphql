package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/r7rainz/go-grpc-graphql/account"
)

type AppConfig struct {
	DatabaseURL string `envconfig:"ACCOUNT_DATABASE_URL" required:"true"`
}

func main() {
	var cfg AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	repository, err := account.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer repository.Close()

	service := account.NewService(repository)
	_ = service

	log.Println("account service connected to database")
}
