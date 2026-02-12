package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/vcnt72/go-boilerplate/internal/domain"
)

type UserRepository struct {
	db sqlx.ExtContext
}

func (t UserRepository) Create(ctx context.Context, spec domain.User) (*domain.User, error) {
	var id int64
	err := t.db.QueryRowxContext(ctx, "INSERT INTO users(name) VALUES($1) RETURNING id", spec.Name).Scan(&id)
	spec.ID = id

	return &spec, err
}

func (t *UserRepository) WithTx(tx sqlx.ExtContext) *UserRepository {
	return &UserRepository{
		db: tx,
	}
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}
