package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ximura/gowallet/internal/core/domain"
	"github.com/ximura/gowallet/internal/repository"
	"github.com/ximura/gowallet/internal/repository/jet/model"
	"gotest.tools/v3/assert"
)

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	db, mock := newMock()
	repo := repository.NewWalletRepo(db)
	defer func() {
		repo.Close()
	}()
	account := uuid.New()
	wallet := model.Wallet{
		ID:       1,
		Account:  account,
		Amount:   100,
		Currency: "usd",
	}

	tests := map[string]struct {
		err   error
		mocks func(m *sqlmock.ExpectedQuery)
	}{
		"Ok": {
			err: nil,
			mocks: func(m *sqlmock.ExpectedQuery) {
				rows := sqlmock.NewRows([]string{"wallet.id", "wallet.account", "wallet.amount", "wallet.currency"}).
					AddRow(wallet.ID, wallet.Account, wallet.Amount, wallet.Currency)
				m.WillReturnRows(rows)
			},
		},
		"ErrNoRows": {
			err: sql.ErrNoRows,
			mocks: func(m *sqlmock.ExpectedQuery) {
				m.WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			query := `SELECT wallet.id AS "wallet.id", wallet.account AS "wallet.account", wallet.amount AS "wallet.amount", wallet.currency AS "wallet.currency" 
			  FROM public.wallet WHERE wallet.id = \$1;`

			q := mock.ExpectQuery(query).WithArgs(wallet.ID)
			tt.mocks(q)

			w, err := repo.Get(ctx, int(wallet.ID))
			if err != nil {
				assert.ErrorContains(t, err, tt.err.Error())
			} else {
				assert.NilError(t, err)
				assert.Equal(t, int32(w.ID), wallet.ID)
				assert.Equal(t, w.Account, wallet.Account)
				assert.Equal(t, int32(w.Amount), wallet.Amount)
				assert.Equal(t, string(w.Currency), wallet.Currency)
			}
		})
	}
}

func TestList(t *testing.T) {
	ctx := context.Background()
	db, mock := newMock()
	repo := repository.NewWalletRepo(db)
	defer func() {
		repo.Close()
	}()
	account := uuid.New()
	wallet := []model.Wallet{{
		ID:       1,
		Account:  account,
		Amount:   100,
		Currency: "usd",
	},
		{
			ID:       2,
			Account:  account,
			Amount:   200,
			Currency: "eur",
		},
	}

	tests := map[string]struct {
		err   error
		mocks func(m *sqlmock.ExpectedQuery)
	}{
		"Ok": {
			err: nil,
			mocks: func(m *sqlmock.ExpectedQuery) {
				rows := sqlmock.NewRows([]string{"wallet.id", "wallet.account", "wallet.amount", "wallet.currency"})
				for _, w := range wallet {
					rows.AddRow(w.ID, w.Account, w.Amount, w.Currency)
				}
				m.WillReturnRows(rows)
			},
		},
		"ErrNoRows": {
			err: sql.ErrNoRows,
			mocks: func(m *sqlmock.ExpectedQuery) {
				m.WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			query := `SELECT wallet.id AS "wallet.id", wallet.account AS "wallet.account", wallet.amount AS "wallet.amount", wallet.currency AS "wallet.currency" 
	FROM public.wallet WHERE wallet.account = \$1;`
			q := mock.ExpectQuery(query).WithArgs(account)
			tt.mocks(q)

			result, err := repo.List(ctx, account)

			if err != nil {
				assert.ErrorContains(t, err, tt.err.Error())
			} else {
				assert.NilError(t, err)
				assert.Equal(t, len(result), len(wallet))
				for i := range wallet {
					assert.Equal(t, int32(result[i].ID), wallet[i].ID)
					assert.Equal(t, result[i].Account, wallet[i].Account)
					assert.Equal(t, int32(result[i].Amount), wallet[i].Amount)
					assert.Equal(t, string(result[i].Currency), wallet[i].Currency)
				}
			}
		})
	}
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	db, mock := newMock()
	repo := repository.NewWalletRepo(db)
	defer func() {
		repo.Close()
	}()
	account := uuid.New()
	wallet := model.Wallet{
		ID:       1,
		Account:  account,
		Amount:   100,
		Currency: "usd",
	}

	tests := map[string]struct {
		err   error
		mocks func(m *sqlmock.ExpectedQuery)
	}{
		"Ok": {
			err: nil,
			mocks: func(m *sqlmock.ExpectedQuery) {
				rows := sqlmock.NewRows([]string{"wallet.id", "wallet.account", "wallet.amount", "wallet.currency"}).
					AddRow(wallet.ID, wallet.Account, wallet.Amount, wallet.Currency)
				m.WillReturnRows(rows)
			},
		},
		"ErrNoRows": {
			err: sql.ErrNoRows,
			mocks: func(m *sqlmock.ExpectedQuery) {
				m.WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {

			query := `INSERT INTO public.wallet \(account, currency\)
				VALUES \(\$1, \$2\)
				RETURNING wallet.id AS "wallet.id", wallet.account AS "wallet.account",
				wallet.amount AS "wallet.amount", wallet.currency AS "wallet.currency";`

			q := mock.ExpectQuery(query).WithArgs(account, wallet.Currency)
			tt.mocks(q)

			result, err := repo.Create(ctx, account, domain.Currency(wallet.Currency))
			if err != nil {
				assert.ErrorContains(t, err, tt.err.Error())
			} else {
				assert.NilError(t, err)
				assert.Equal(t, int32(result.ID), wallet.ID)
				assert.Equal(t, result.Account, wallet.Account)
				assert.Equal(t, int32(result.Amount), wallet.Amount)
				assert.Equal(t, string(result.Currency), wallet.Currency)
			}
		})
	}
}

func TestHasTransaction(t *testing.T) {
	ctx := context.Background()
	db, mock := newMock()
	repo := repository.NewWalletRepo(db)
	defer func() {
		repo.Close()
	}()
	transaction := domain.Transaction{
		ID:       uuid.New(),
		WalletID: 1,
		Amount:   10,
		Currency: "usd",
	}

	tests := map[string]struct {
		result bool
		err    error
		mocks  func(m *sqlmock.ExpectedQuery)
	}{
		"True": {
			result: true,
			err:    nil,
			mocks: func(m *sqlmock.ExpectedQuery) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				m.WillReturnRows(rows)
			},
		},
		"False": {
			result: false,
			err:    nil,
			mocks: func(m *sqlmock.ExpectedQuery) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				m.WillReturnRows(rows)
			},
		},
		"ErrNoRows": {
			err: sql.ErrNoRows,
			mocks: func(m *sqlmock.ExpectedQuery) {
				m.WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			query := `SELECT COUNT\(transaction.wallet_id\) AS "count" FROM public.transaction WHERE \(transaction.wallet_id = \$1\) AND \(transaction.transaction_id = \$2\);`

			q := mock.ExpectQuery(query).WithArgs(transaction.WalletID, transaction.ID)
			tt.mocks(q)

			ok, err := repo.HasTransaction(ctx, transaction)
			if tt.err != nil {
				assert.ErrorContains(t, err, tt.err.Error())
			} else {
				assert.NilError(t, err)
				assert.Equal(t, ok, tt.result)
			}
		})
	}
}

func TestProcessTransaction(t *testing.T) {
	ctx := context.Background()
	transaction := domain.Transaction{
		ID:       uuid.New(),
		WalletID: 1,
		Amount:   10,
		Currency: "usd",
	}
	wallet := model.Wallet{
		ID:       1,
		Account:  uuid.New(),
		Amount:   100,
		Currency: "usd",
	}

	tests := map[string]struct {
		err   error
		mocks func(m sqlmock.Sqlmock, err error)
	}{
		"Transaction_Exist": {
			err: errors.New("pq: duplicate key value violates unique constraint \"transaction_pkey\""),
			mocks: func(mock sqlmock.Sqlmock, err error) {
				mock.MatchExpectationsInOrder(true)
				mock.ExpectBegin()

				query := `INSERT INTO public.transaction \(wallet_id, transaction_id\) VALUES \(\$1, \$2\);`
				mock.ExpectExec(query).WithArgs(transaction.WalletID, transaction.ID).WillReturnError(err)
				mock.ExpectRollback()
			},
		},
		"negative": {
			err: errors.New("pq: new row for relation \"wallet\" violates check constraint \"positive_amount\""),
			mocks: func(mock sqlmock.Sqlmock, err error) {
				mock.MatchExpectationsInOrder(true)
				mock.ExpectBegin()

				query := `INSERT INTO public.transaction \(wallet_id, transaction_id\) VALUES \(\$1, \$2\);`
				mock.ExpectExec(query).WithArgs(transaction.WalletID, transaction.ID).WillReturnResult(sqlmock.NewResult(0, 0))

				query = `UPDATE public.wallet
					SET amount = \(wallet.amount \+ \$1\)
					WHERE \(wallet.id = \$2\) AND \(wallet.currency = \$3::text\)
					RETURNING wallet.id AS "wallet.id", wallet.account AS "wallet.account", wallet.amount AS "wallet.amount", wallet.currency AS "wallet.currency";`
				mock.ExpectQuery(query).WithArgs(transaction.Amount, transaction.WalletID, transaction.Currency).WillReturnError(err)
				mock.ExpectRollback()
			},
		},
		"Ok": {
			err: nil,
			mocks: func(mock sqlmock.Sqlmock, err error) {
				mock.MatchExpectationsInOrder(true)
				mock.ExpectBegin()

				query := `INSERT INTO public.transaction \(wallet_id, transaction_id\) VALUES \(\$1, \$2\);`
				mock.ExpectExec(query).WithArgs(transaction.WalletID, transaction.ID).WillReturnResult(sqlmock.NewResult(0, 0))

				rows := sqlmock.NewRows([]string{"wallet.id", "wallet.account", "wallet.amount", "wallet.currency"}).
					AddRow(wallet.ID, wallet.Account, wallet.Amount, wallet.Currency)

				query = `UPDATE public.wallet
					SET amount = \(wallet.amount \+ \$1\)
					WHERE \(wallet.id = \$2\) AND \(wallet.currency = \$3::text\)
					RETURNING wallet.id AS "wallet.id", wallet.account AS "wallet.account", wallet.amount AS "wallet.amount", wallet.currency AS "wallet.currency";`
				mock.ExpectQuery(query).WithArgs(transaction.Amount, transaction.WalletID, transaction.Currency).WillReturnRows(rows)
				mock.ExpectCommit()
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			db, mock := newMock()
			repo := repository.NewWalletRepo(db)
			defer func() {
				repo.Close()
			}()

			tt.mocks(mock, tt.err)

			_, err := repo.ProcessTransaction(ctx, transaction)
			if tt.err != nil {
				assert.ErrorContains(t, err, tt.err.Error())
			} else {
				assert.NilError(t, err)
			}
		})
	}
}
