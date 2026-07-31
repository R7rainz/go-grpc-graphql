package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/r7rainz/go-grpc-graphql/catalog/internal/domain"
)

type Service interface {
	PostProduct(ctx context.Context, name string, description string, price float64) (*domain.Product, error)
	GetProduct(ctx context.Context, id string) (*domain.Product, error)
	GetProducts(ctx context.Context, skip int, take int) ([]domain.Product, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]domain.Product, error)
	SearchProducts(ctx context.Context, query string, skip int, take int) ([]domain.Product, error)
}

type Repository interface {
	PutProduct(ctx context.Context, p domain.Product) error
	GetProductByID(ctx context.Context, id string) (*domain.Product, error)
	ListProducts(ctx context.Context, skip int, take int) ([]domain.Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]domain.Product, error)
	SearchProducts(ctx context.Context, query string, skip int, take int) ([]domain.Product, error)
}

type catalogService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &catalogService{repository: r}
}

// paging clamps caller-supplied paging so a bad request cannot ask for the
// whole index at once.
func paging(skip int, take int) (int, int) {
	if skip < 0 {
		skip = 0
	}
	if take <= 0 || take > 100 {
		take = 100
	}
	return skip, take
}

func (s *catalogService) PostProduct(ctx context.Context, name string, description string, price float64) (*domain.Product, error) {
	product := domain.Product{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		Price:       price,
	}

	if err := s.repository.PutProduct(ctx, product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *catalogService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	return s.repository.GetProductByID(ctx, id)
}

func (s *catalogService) GetProducts(ctx context.Context, skip int, take int) ([]domain.Product, error) {
	skip, take = paging(skip, take)
	return s.repository.ListProducts(ctx, skip, take)
}

func (s *catalogService) GetProductsByIDs(ctx context.Context, ids []string) ([]domain.Product, error) {
	if len(ids) == 0 {
		return []domain.Product{}, nil
	}
	return s.repository.ListProductsWithIDs(ctx, ids)
}

func (s *catalogService) SearchProducts(ctx context.Context, query string, skip int, take int) ([]domain.Product, error) {
	skip, take = paging(skip, take)
	return s.repository.SearchProducts(ctx, query, skip, take)
}
