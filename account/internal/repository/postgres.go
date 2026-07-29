package repository

import (
	"context"
	"database/sql"

	"github.com/r7rainz/go-grpc-graphql/account/internal/domain"

	_ "github.com/lib/pq"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*postgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	// pinging the database just to see if the connection is established or not
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return &postgresRepository{db: db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgresRepository) PutAccount(ctx context.Context, a domain.Account) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO account(id, name) VALUES($1, $2)", a.ID, a.Name)
	return err
}

func (r *postgresRepository) GetAccountByID(ctx context.Context, id string) (*domain.Account, error) {
	rows := r.db.QueryRowContext(ctx, "SELECT id, name FROM account WHERE id = $1", id)
	a := &domain.Account{}
	if err := rows.Scan(&a.ID, &a.Name); err != nil {
		return nil, err
	}
	return a, nil
}

func (r *postgresRepository) ListAccounts(ctx context.Context, skip int, take int) ([]domain.Account, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM account ORDER BY id DESC OFFSET $1 LIMIT $2", skip, take)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := []domain.Account{}

	for rows.Next() {
		a := &domain.Account{}
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, err
		}
		accounts = append(accounts, *a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
