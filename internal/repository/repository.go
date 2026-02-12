// Package repository
package repository

import "github.com/jmoiron/sqlx"

type Repositories struct {
	UserRepository   *UserRepository
	WalletRepository *WalletRepository
	LedgerRepository *LedgerRepository
	TxProvider       *TxProvider
}

func New(db *sqlx.DB) Repositories {
	return Repositories{
		UserRepository:   NewUserRepository(db),
		WalletRepository: NewWalletRepository(db),
		LedgerRepository: NewLedgerRepository(db),
		TxProvider:       NewTxProvider(db),
	}
}
