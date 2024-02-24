package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ximura/gowallet/internal/core/domain"
)

//go:generate  go run github.com/golang/mock/mockgen@v1.6.0 -source=$GOFILE -package=mocks -destination=mocks/repository_mock.go
type WalletRepository interface {
	// Creates new wallet for account
	Create(context.Context, uuid.UUID, domain.Currency) (domain.Wallet, error)
	// Return  list of wallets linked to account
	List(context.Context, uuid.UUID) ([]domain.Wallet, error)
	// Return current state of account wallet
	Get(context.Context, int) (domain.Wallet, error)
	// Check if a transaction with the same id was already processed
	HasTransaction(context.Context, domain.Transaction) (bool, error)
	// Execute transaction for account wallet
	ProcessTransaction(context.Context, domain.Transaction) (domain.Wallet, error)
}
