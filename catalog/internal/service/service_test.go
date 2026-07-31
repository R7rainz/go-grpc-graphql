package service

import (
	"context"
	"testing"

	"github.com/r7rainz/go-grpc-graphql/catalog/internal/domain"
)

// fakeRepository records what the service asked for, so the tests can assert on
// it without a running Elasticsearch.
type fakeRepository struct {
	put       domain.Product
	skip      int
	take      int
	Repository // unimplemented methods panic if a test ever reaches them
}

func (f *fakeRepository) PutProduct(_ context.Context, p domain.Product) error {
	f.put = p
	return nil
}

func (f *fakeRepository) ListProducts(_ context.Context, skip int, take int) ([]domain.Product, error) {
	f.skip, f.take = skip, take
	return nil, nil
}

func (f *fakeRepository) SearchProducts(_ context.Context, _ string, skip int, take int) ([]domain.Product, error) {
	f.skip, f.take = skip, take
	return nil, nil
}

func TestPostProductGeneratesID(t *testing.T) {
	repo := &fakeRepository{}
	p, err := NewService(repo).PostProduct(context.Background(), "shoe", "blue", 49.99)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID == "" {
		t.Error("expected a generated ID")
	}
	if repo.put != *p {
		t.Errorf("stored %+v, returned %+v", repo.put, *p)
	}
}

func TestPagingIsClamped(t *testing.T) {
	for _, tc := range []struct {
		skip, take         int
		wantSkip, wantTake int
	}{
		{0, 10, 0, 10},
		{-5, 10, 0, 10},
		{20, 0, 20, 100},
		{20, 1000, 20, 100},
	} {
		repo := &fakeRepository{}
		s := NewService(repo)

		if _, err := s.GetProducts(context.Background(), tc.skip, tc.take); err != nil {
			t.Fatal(err)
		}
		if repo.skip != tc.wantSkip || repo.take != tc.wantTake {
			t.Errorf("GetProducts(%d, %d) = (%d, %d), want (%d, %d)",
				tc.skip, tc.take, repo.skip, repo.take, tc.wantSkip, tc.wantTake)
		}

		if _, err := s.SearchProducts(context.Background(), "q", tc.skip, tc.take); err != nil {
			t.Fatal(err)
		}
		if repo.skip != tc.wantSkip || repo.take != tc.wantTake {
			t.Errorf("SearchProducts(%d, %d) = (%d, %d), want (%d, %d)",
				tc.skip, tc.take, repo.skip, repo.take, tc.wantSkip, tc.wantTake)
		}
	}
}
