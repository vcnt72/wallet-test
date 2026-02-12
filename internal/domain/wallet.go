package domain

import (
	"errors"
	"time"
)

var (
	ErrWalletNotFound   = errors.New("error wallet not found")
	ErrInsufficientFund = errors.New("error insuficient fund")
	ErrInvalidAmount    = errors.New("error invalid amount")
)

type Wallet struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Balance   int64     `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
