package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/vcnt72/go-boilerplate/internal/domain"
)

type WalletRepository struct {
	db sqlx.ExtContext
}

func (w WalletRepository) Create(ctx context.Context, spec domain.Wallet) (*domain.Wallet, error) {
	var id int64
	err := w.db.QueryRowxContext(ctx, "INSERT INTO wallets(user_id, balance) VALUES($1,$2) RETURNING id", spec.UserID, spec.Balance).Scan(&id)
	spec.ID = id
	return &spec, err
}

func (w WalletRepository) GetByUserID(ctx context.Context, userID int64) (*domain.Wallet, error) {
	var wallet domain.Wallet

	err := w.db.QueryRowxContext(ctx, "SELECT id, balance, user_id, created_at, updated_at FROM wallets WHERE user_id = $1", userID).
		StructScan(&wallet)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrWalletNotFound
		}

		return nil, err
	}

	return &wallet, nil
}

func (w WalletRepository) DecreaseBalance(ctx context.Context, amount, userID int64) (int64, error) {
	if amount <= 0 {
		return 0, domain.ErrInvalidAmount
	}

	var b int64
	err := w.db.QueryRowxContext(ctx,
		"UPDATE wallets SET balance = balance - $1, updated_at = now() WHERE user_id = $2 AND balance >= $1 RETURNING balance", amount, userID).
		Scan(&b)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrInsufficientFund
		}

		return 0, err
	}

	return b, nil
}

func (w WalletRepository) WithTx(tx sqlx.ExtContext) *WalletRepository {
	return &WalletRepository{
		db: tx,
	}
}

func NewWalletRepository(db *sqlx.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}
