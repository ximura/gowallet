package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ximura/gowallet/internal/core/domain"
	"github.com/ximura/gowallet/internal/core/ports"
)

var _ ports.WalletService = (*WalletService)(nil)

var ErrInvalitTransactionAmount = errors.New("invalid transaction amount")
var ErrDuplicateTransaction = errors.New("duplicate transaction")

type WalletService struct {
	repo ports.WalletRepository
}

func NewWalletService(repo ports.WalletRepository) WalletService {
	return WalletService{repo: repo}
}

func (w *WalletService) Create(ctx context.Context, account uuid.UUID, currency domain.Currency) (domain.Wallet, error) {
	return w.repo.Create(ctx, account, currency)
}

func (w *WalletService) Get(ctx context.Context, id int) (domain.Wallet, error) {
	return w.repo.Get(ctx, id)
}

func (w *WalletService) List(ctx context.Context, account uuid.UUID) ([]domain.Wallet, error) {
	return w.repo.List(ctx, account)
}

func (w *WalletService) ProcessTransaction(ctx context.Context, transaction domain.Transaction) (domain.Wallet, error) {
	ok, err := w.repo.HasTransaction(ctx, transaction)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("can't get transaction: %w", err)
	}
	if ok {
		return domain.Wallet{}, ErrDuplicateTransaction
	}

	wallet, err := w.repo.Get(ctx, transaction.WalletID)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("can't get wallet %d: %w", transaction.WalletID, err)
	}

	if wallet.Currency != transaction.Currency {
		return domain.Wallet{}, fmt.Errorf("wallet currency different from transaction, %s != %s", wallet.Currency, transaction.Currency)
	}

	nAmount := wallet.Amount + transaction.Amount
	if nAmount < 0 {
		return domain.Wallet{}, ErrInvalitTransactionAmount
	}

	return w.repo.ProcessTransaction(ctx, transaction)
}
