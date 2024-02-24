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

func TestProcessTransaction(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	tErr := errors.New("test error")
	transaction := domain.Transaction{
		ID:       uuid.New(),
		WalletID: 1,
		Amount:   100,
		Currency: "usd",
	}

	tests := map[string]struct {
		err   error
		mocks func(m *mocks.MockWalletRepository)
	}{
		"can't get transaction": {
			err: fmt.Errorf("can't get transaction: %w", tErr),
			mocks: func(m *mocks.MockWalletRepository) {
				m.EXPECT().HasTransaction(ctx, transaction).Return(true, tErr)
			},
		},
		"ErrDuplicateTransaction": {
			err: service.ErrDuplicateTransaction,
			mocks: func(m *mocks.MockWalletRepository) {
				m.EXPECT().HasTransaction(ctx, transaction).Return(true, nil)
			},
		},
		"can't get wallet for": {
			err: fmt.Errorf("can't get wallet %d: %w", transaction.WalletID, tErr),
			mocks: func(m *mocks.MockWalletRepository) {
				m.EXPECT().HasTransaction(ctx, transaction).Return(false, nil)
				m.EXPECT().Get(ctx, transaction.WalletID).Return(domain.Wallet{}, tErr)
			},
		},
		"wallet currency different from transaction": {
			err: fmt.Errorf("wallet currency different from transaction, %s != %s", "eur", transaction.Currency),
			mocks: func(m *mocks.MockWalletRepository) {
				m.EXPECT().HasTransaction(ctx, transaction).Return(false, nil)
				m.EXPECT().Get(ctx, transaction.WalletID).Return(domain.Wallet{
					Currency: "eur",
				}, nil)
			},
		},
		"ErrInvalitTransactionAmount": {
			err: service.ErrInvalitTransactionAmount,
			mocks: func(m *mocks.MockWalletRepository) {
				m.EXPECT().HasTransaction(ctx, transaction).Return(false, nil)
				m.EXPECT().Get(ctx, transaction.WalletID).Return(domain.Wallet{
					Currency: "usd",
					Amount:   -1 * (transaction.Amount + 10),
				}, nil)
			},
		},
		"Transaction error": {
			err: tErr,
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
			err: nil,
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
