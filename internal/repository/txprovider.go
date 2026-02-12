package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TxProvider struct {
	db *sqlx.DB
}

func (t TxProvider) Tx(ctx context.Context, txFunc func(sqlx.ExtContext) error) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	err = txFunc(tx)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	return tx.Commit()
}

func NewTxProvider(db *sqlx.DB) *TxProvider {
	return &TxProvider{db: db}
}
