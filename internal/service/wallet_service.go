package service

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/vcnt72/go-boilerplate/internal/domain"
	"github.com/vcnt72/go-boilerplate/internal/repository"
)

type WalletService struct {
	walletRepository *repository.WalletRepository
	ledgerRepository *repository.LedgerRepository
	txProvider       *repository.TxProvider
}

func (w WalletService) GetByUserID(ctx context.Context, userID int64) (*domain.Wallet, error) {
	wallet, err := w.walletRepository.GetByUserID(ctx, userID)

	return wallet, err
}

type WithdrawWalletSpec struct {
	UserID         int64
	IdempotencyKey string
	Amount         int64
}

type WithdrawalResult struct {
	UserID  int64
	Balance int64
	Amount  int64
}

func (w WalletService) Withdraw(ctx context.Context, spec WithdrawWalletSpec) (*WithdrawalResult, error) {
	var balance int64
	var appErr error
	err := w.txProvider.Tx(ctx, func(tx sqlx.ExtContext) error {
		wallet, err := w.walletRepository.WithTx(tx).GetByUserID(ctx, spec.UserID)
		if err != nil {
			return err
		}

		ledger, err := w.ledgerRepository.WithTx(tx).Create(ctx, domain.Ledger{
			IdempotencyKey: spec.IdempotencyKey,
			Type:           domain.LedgerTypeWithdraw,
			WalletID:       wallet.ID,
			Status:         domain.LedgerStatusProcessing,
			Amount:         spec.Amount,
		})
		if err != nil {
			return err
		}

		balance, err = w.walletRepository.WithTx(tx).DecreaseBalance(ctx, spec.Amount, spec.UserID)
		if err != nil {
			if errors.Is(err, domain.ErrInsufficientFund) {
				errCode := "INSUFFICIENT_FUND"
				ledger.Status = domain.LedgerStatusFailed
				ledger.ErrorCode = &errCode
				ledger.ResultBalance = &wallet.Balance
				uerr := w.ledgerRepository.WithTx(tx).Update(ctx, *ledger)
				appErr = err
				return uerr
			}

			return err
		}

		wallet.Balance = balance
		ledger.Status = domain.LedgerStatusSucceed
		ledger.ResultBalance = &balance
		if err := w.ledgerRepository.WithTx(tx).Update(ctx, *ledger); err != nil {
			return err
		}

		return nil
	})

	if errors.Is(err, domain.ErrLedgerConflict) {
		withdrawResult, err2 := w.handleConflictWithdraw(ctx, spec)

		if err2 != nil {
			return nil, errors.Join(err, err2)
		}

		return withdrawResult, nil
	}

	if err != nil {
		return nil, err
	}

	if appErr != nil {
		return nil, appErr
	}

	return &WithdrawalResult{
		UserID:  spec.UserID,
		Amount:  spec.Amount,
		Balance: balance,
	}, err
}

func (w WalletService) handleConflictWithdraw(ctx context.Context, spec WithdrawWalletSpec) (*WithdrawalResult, error) {
	l, err := w.ledgerRepository.GetByIdempotencyKey(ctx, spec.IdempotencyKey)
	if err != nil {
		return nil, err
	}

	if l.Amount != spec.Amount || l.Type != domain.LedgerTypeWithdraw {
		return nil, domain.ErrIdempotencyKeyReused
	}

	switch l.Status {
	case domain.LedgerStatusSucceed:
		if l.ResultBalance == nil {
			return nil, domain.ErrRequestInProgress
		}
		return &WithdrawalResult{
			UserID:  spec.UserID,
			Amount:  spec.Amount,
			Balance: *l.ResultBalance,
		}, nil

	case domain.LedgerStatusFailed:
		if l.ErrorCode != nil && *l.ErrorCode == "INSUFFICIENT_FUND" {
			return nil, domain.ErrInsufficientFund
		}
		return nil, domain.ErrWithdrawFailed

	default:
		return nil, domain.ErrRequestInProgress
	}
}

func NewWalletService(walletRepository *repository.WalletRepository, ledgerRepository *repository.LedgerRepository, txProvider *repository.TxProvider) *WalletService {
	return &WalletService{
		walletRepository: walletRepository,
		ledgerRepository: ledgerRepository,
		txProvider:       txProvider,
	}
}
