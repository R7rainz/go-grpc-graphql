package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/r7rainz/go-grpc-graphql/account/internal/domain"
)

type Service interface {
	PostAccount(ctx context.Context, name string) (*domain.Account, error)
	GetAccount(ctx context.Context, id string) (*domain.Account, error)
	GetAccounts(ctx context.Context, skip int, take int) ([]domain.Account, error)
}

type Repository interface {
	PutAccount(ctx context.Context, a domain.Account) error
	GetAccountByID(ctx context.Context, id string) (*domain.Account, error)
	ListAccounts(ctx context.Context, skip int, take int) ([]domain.Account, error)
}

type accountService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &accountService{repository: r}
}

func (s *accountService) PostAccount(ctx context.Context, name string) (*domain.Account, error) {
	account := domain.Account{
		ID:   uuid.NewString(),
		Name: name,
	}

	if err := s.repository.PutAccount(ctx, account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *accountService) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	return s.repository.GetAccountByID(ctx, id)
}

func (s *accountService) GetAccounts(ctx context.Context, skip int, take int) ([]domain.Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListAccounts(ctx, skip, take)
}
