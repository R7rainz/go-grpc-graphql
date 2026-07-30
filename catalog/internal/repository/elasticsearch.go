package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/olivere/elastic"
	"github.com/r7rainz/go-grpc-graphql/catalog/internal/domain"
)

var ErrNotFound = errors.New("Entity not found")

type elasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desciption"`
	Price       string `json:"price"`
}

func NewElasticRepository(URL string) (*elasticRepository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(URL),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}

	return &elasticRepository{
		client,
	}, nil
}

func (r *elasticRepository) Close() {
}

func (r *elasticRepository) PutProduct(ctx context.Context, p domain.Product) error {
	_, err := r.client.Index().Index("catalog").Type("product").Id(p.ID).BodyJson(productDocument{
		Name: p.Name, Description: p.Description, Price: p.Price,
	}).Do(ctx)

	return err
}

func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	res, err := r.client.Get().Index("catalog").Type("product").Id(id).Do(ctx)
	if err != nil {
		return nil, err
	}
	if !res.Found {
		return nil, ErrNotFound
	}
	p := productDocument{}
	if err = json.Unmarshal(*res.Source, &p); err != nil {
		return nil, err
	}

	return &domain.Product{
		ID:          id,
		Name:        p.Name,
		Description: p.Price,
		Price:       p.Price,
	}, nil
}

func (r *elasticRepository) ListProducts(ctx context.Context, skip int, take int) ([]domain.Product, error) {
	res, err := r.client.Search().Index("catalog").Type("product").Query(elastic.NewMatchAllQuery()).From(int(skip)).Size(int(take)).Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// hits are all the  values matched
	products := []domain.Product{}
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err = json.Unmarshal(*hit.Source, &p); err == nil {
			products = append(products, domain.Product{
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, err
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]domain.Product, error) {
	items := []*elastic.MultiGetItem{}
	for _, id := range ids {
		items = append(items, elastic.NewMultiGetItem().Index("catalog").Type("product").Id(id))
	}

	res, err := r.client.MultiGet().Add(items...).Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	products := []domain.Product{}
	for _, doc := range res.Docs {
		p := productDocument{}
		if err := json.Unmarshal(*doc.Source, &p); err != nil {
			return nil, err
		}

		products = append(products, domain.Product{
			ID:          doc.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return products, err
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip int, take int) ([]domain.Product, error) {
	res, err := r.client.Search().Index("catalog").Type("product").Query(elastic.NewMultiMatchQuery(query, "name", "description")).From(skip).Size(take).Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []domain.Product{}
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err := json.Unmarshal(*hit.Source, &p); err != nil {
			return nil, err
		}

		products = append(products, domain.Product{
			ID:          hit.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return products, err
}
