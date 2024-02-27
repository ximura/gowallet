package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ximura/gowallet/internal/core/domain"
	"github.com/ximura/gowallet/internal/core/ports/mocks"
	"github.com/ximura/gowallet/internal/core/service"
	"gotest.tools/v3/assert"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	account := uuid.New()

	tests := map[string]struct {
		currency domain.Currency
		err      error
		mocks    func(c domain.Currency, m *mocks.MockWalletRepository)
	}{
		"usd": {
			currency: "usd",
			err:      nil,
			mocks: func(c domain.Currency, m *mocks.MockWalletRepository) {
				m.EXPECT().Create(ctx, account, c).Return(domain.Wallet{}, nil)
			},
		},
		"eur": {
			currency: "EUR",
			err:      nil,
			mocks: func(c domain.Currency, m *mocks.MockWalletRepository) {
				m.EXPECT().Create(ctx, account, c).Return(domain.Wallet{}, nil)
			},
		},
		"jpy": {
			currency: "jPy",
			err:      nil,
			mocks: func(c domain.Currency, m *mocks.MockWalletRepository) {
				m.EXPECT().Create(ctx, account, c).Return(domain.Wallet{}, nil)
			},
		},
		"error": {
			currency: "test",
			err:      service.ErrUnsuportedCurrency,
			mocks: func(c domain.Currency, m *mocks.MockWalletRepository) {
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			repository := mocks.NewMockWalletRepository(ctrl)
			wallet := service.NewWalletService(repository)
			tt.mocks(tt.currency, repository)
			_, err := wallet.Create(ctx, account, tt.currency)
			if tt.err == nil {
				assert.NilError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err.Error())
			}
		})
	}
}

func TestProcessTransaction(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	tErr := errors.New("test error")
	transaction := domain.Transaction{
		ID:       uuid.New(),
		WalletID: 1,
		Amount:   100,
	}

	tests := map[string]struct {
		currency domain.Currency
		err      error
		mocks    func(m *mocks.MockWalletRepository)
	}{
		"can't get transaction": {
			currency: "usd",
			err:      fmt.Errorf("can't get transaction: %w", tErr),
			mocks: func(m *mocks.MockWalletRepository) {
				m.EXPECT().HasTransaction(ctx, transaction).Return(true, tErr)
			},
		},
		"ErrUnsuportedCurrency": {
			currency: "test",
			err:      service.ErrUnsuportedCurrency,
			mocks: func(m *mocks.MockWalletRepository) {
			},
		},
		"ErrDuplicateTransaction": {
			currency: "usd",
			err:      service.ErrDuplicateTransaction,
			mocks: func(m *mocks.MockWalletRepository) {
				m.EXPECT().HasTransaction(ctx, transaction).Return(true, nil)
			},
		},
		"can't get wallet for": {
			currency: "usd",
			err:      fmt.Errorf("can't get wallet %d: %w", transaction.WalletID, tErr),
			mocks: func(m *mocks.MockWalletRepository) {
				m.EXPECT().HasTransaction(ctx, transaction).Return(false, nil)
				m.EXPECT().Get(ctx, transaction.WalletID).Return(domain.Wallet{}, tErr)
			},
		},
		"wallet currency different from transaction": {
			currency: "usd",
			err:      fmt.Errorf("wallet currency different from transaction, %s != %s", "eur", transaction.Currency),
			mocks: func(m *mocks.MockWalletRepository) {
				m.EXPECT().HasTransaction(ctx, transaction).Return(false, nil)
				m.EXPECT().Get(ctx, transaction.WalletID).Return(domain.Wallet{
					Currency: "eur",
				}, nil)
			},
		},
		"ErrInvalitTransactionAmount": {
			currency: "usd",
			err:      service.ErrInvalitTransactionAmount,
			mocks: func(m *mocks.MockWalletRepository) {
				m.EXPECT().HasTransaction(ctx, transaction).Return(false, nil)
				m.EXPECT().Get(ctx, transaction.WalletID).Return(domain.Wallet{
					Currency: "usd",
					Amount:   -1 * (transaction.Amount + 10),
				}, nil)
			},
		},
		"Transaction error": {
			currency: "usd",
			err:      tErr,
			mocks: func(m *mocks.MockWalletRepository) {
				m.EXPECT().HasTransaction(ctx, transaction).Return(false, nil)
				m.EXPECT().Get(ctx, transaction.WalletID).Return(domain.Wallet{
					Currency: "usd",
					Amount:   transaction.Amount + 10,
				}, nil)
				m.EXPECT().ProcessTransaction(ctx, transaction).Return(domain.Wallet{}, tErr)
			},
		},
		"Ok": {
			currency: "usd",
			err:      nil,
			mocks: func(m *mocks.MockWalletRepository) {
				m.EXPECT().HasTransaction(ctx, transaction).Return(false, nil)
				m.EXPECT().Get(ctx, transaction.WalletID).Return(domain.Wallet{
					Currency: "usd",
					Amount:   transaction.Amount + 10,
				}, nil)
				m.EXPECT().ProcessTransaction(ctx, transaction).Return(domain.Wallet{}, nil)
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			repository := mocks.NewMockWalletRepository(ctrl)
			transaction.Currency = tt.currency
			tt.mocks(repository)
			wallet := service.NewWalletService(repository)

			_, err := wallet.ProcessTransaction(ctx, transaction)
			if tt.err == nil {
				assert.NilError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err.Error())
			}
		})
	}
}
