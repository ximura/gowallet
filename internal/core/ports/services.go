package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ximura/gowallet/internal/core/domain"
)

type WalletService interface {
	// Creates new wallet for account
	Create(context.Context, uuid.UUID, domain.Currency) (domain.Wallet, error)
	// Return current state of account wallet
	Get(context.Context, int) (domain.Wallet, error)
	// Return  list of wallets linked to account
	List(context.Context, uuid.UUID) ([]domain.Wallet, error)
	// Execute transaction for account wallet
	ProcessTransaction(context.Context, domain.Transaction) (domain.Wallet, error)
}
