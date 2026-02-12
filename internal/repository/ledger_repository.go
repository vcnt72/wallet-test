package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/vcnt72/go-boilerplate/internal/domain"
)

type LedgerRepository struct {
	db sqlx.ExtContext
}

func (l LedgerRepository) Create(ctx context.Context, ledger domain.Ledger) (*domain.Ledger, error) {
	var id int64
	var status string
	err := l.db.QueryRowxContext(ctx, "INSERT INTO ledgers(idempotency_key, amount, type, status, wallet_id) VALUES($1,$2,$3,$4,$5) ON CONFLICT DO NOTHING RETURNING id, status",
		ledger.IdempotencyKey,
		ledger.Amount,
		ledger.Type,
		ledger.Status,
		ledger.WalletID,
	).Scan(&id, &status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrLedgerConflict
		}

		return nil, err
	}

	ledger.ID = id

	return &ledger, nil
}

func (l LedgerRepository) Update(ctx context.Context, spec domain.Ledger) error {
	_, err := l.db.ExecContext(ctx, "UPDATE ledgers SET status = $1, error_code = $2, result_balance = $3, updated_at = now() WHERE id = $4",

		spec.Status,
		spec.ErrorCode,
		spec.ResultBalance,
		spec.ID,
	)

	return err
}

func (l LedgerRepository) GetByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Ledger, error) {
	var ledger domain.Ledger

	err := l.db.QueryRowxContext(ctx, "SELECT id, idempotency_key, wallet_id, type, status, amount, result_balance, error_code, created_at, updated_at FROM ledgers WHERE idempotency_key = $1", idempotencyKey).
		StructScan(&ledger)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrLedgerNotFound
		}

		return nil, err
	}

	return &ledger, nil
}

func (l *LedgerRepository) WithTx(tx sqlx.ExtContext) *LedgerRepository {
	return &LedgerRepository{
		db: tx,
	}
}

func NewLedgerRepository(db sqlx.ExtContext) *LedgerRepository {
	return &LedgerRepository{
		db: db,
	}
}
