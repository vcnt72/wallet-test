// Package service
package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/vcnt72/go-boilerplate/internal/domain"
	"github.com/vcnt72/go-boilerplate/internal/repository"
)

type UserService struct {
	userRepository   *repository.UserRepository
	walletRepository *repository.WalletRepository
	ledgerRepository *repository.LedgerRepository
	txProvider       *repository.TxProvider
}

type CreateUserSpec struct {
	Name    string
	Balance int64
}

func (t UserService) Create(ctx context.Context, spec CreateUserSpec) (*domain.User, error) {
	var userObj *domain.User
	err := t.txProvider.Tx(ctx, func(tx sqlx.ExtContext) error {
		user, err := t.userRepository.WithTx(tx).Create(ctx, domain.User{
			Name: spec.Name,
		})
		if err != nil {
			return errors.Join(errors.New("UserService.Create: error on user repository create"), err)
		}

		userObj = user

		wallet, err := t.walletRepository.WithTx(tx).Create(ctx, domain.Wallet{
			Balance: spec.Balance,
			UserID:  user.ID,
		})
		if err != nil {
			return errors.Join(errors.New("UserService.Create: error on wallet repository create"), err)
		}

		_, err = t.ledgerRepository.WithTx(tx).Create(ctx, domain.Ledger{
			WalletID:       wallet.ID,
			Amount:         spec.Balance,
			IdempotencyKey: uuid.NewString(),
			Type:           domain.LedgerTypeInit,
			Status:         domain.LedgerStatusSucceed,
		})
		if err != nil {
			return errors.Join(errors.New("UserService.Create: error on ledger repository create"), err)
		}

		return nil
	})

	return userObj, err
}

func NewUserService(userRepository *repository.UserRepository, walletRepository *repository.WalletRepository, ledgerRepository *repository.LedgerRepository, txProvider *repository.TxProvider) *UserService {
	return &UserService{
		userRepository:   userRepository,
		walletRepository: walletRepository,
		ledgerRepository: ledgerRepository,
		txProvider:       txProvider,
	}
}
