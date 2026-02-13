// Package domain
package domain

import (
	"errors"
	"time"
)

type LedgerStatus = string

var (
	LedgerStatusSucceed    = "SUCCEED"
	LedgerStatusFailed     = "FAILED"
	LedgerStatusProcessing = "PROCESSING"
)

type LedgerType = string

var (
	LedgerTypeWithdraw = "WITHDRAW"
	LedgerTypeInit     = "INIT"
)

var (
	ErrLedgerConflict       = errors.New("error ledger unique constraint")
	ErrLedgerNotFound       = errors.New("error ledger not found")
	ErrIdempotencyKeyReused = errors.New("error idempotency key reused")
	ErrWithdrawFailed       = errors.New("error withdraw failed")
	ErrRequestInProgress    = errors.New("error request in progress")
)

type Ledger struct {
	ID             int64        `db:"id"`
	IdempotencyKey string       `db:"idempotency_key"`
	Type           LedgerType   `db:"type"`
	Status         LedgerStatus `db:"status"`
	WalletID       int64        `db:"wallet_id"`
	Amount         int64        `db:"amount"`
	ResultBalance  *int64       `db:"result_balance"`
	ErrorCode      *string      `db:"error_code"`
	CreatedAt      time.Time    `db:"created_at"`
	UpdatedAt      time.Time    `db:"updated_at"`
}
