package repository

import (
	"context"
	"database/sql"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ximura/gowallet/internal/core/domain"
	"github.com/ximura/gowallet/internal/core/ports"
	"github.com/ximura/gowallet/internal/repository/jet/table"
)

var _ ports.WalletRepository = (*WalletRepo)(nil)

type WalletRepo struct {
	db          *sql.DB
	wallet      table.WalletTable
	transaction table.TransactionTable
}

func NewWalletRepo(db *sql.DB) WalletRepo {
	return WalletRepo{
		db:          db,
		wallet:      *table.Wallet,
		transaction: *table.Transaction,
	}
}

func (r *WalletRepo) Close() error {
	return r.db.Close()
}

func (r *WalletRepo) Create(ctx context.Context, account uuid.UUID, currency domain.Currency) (domain.Wallet, error) {
	query := r.wallet.INSERT(
		r.wallet.Account,
		r.wallet.Currency,
	).VALUES(account, currency).
		RETURNING(r.wallet.AllColumns.Except(r.wallet.CreatedAt, r.wallet.UpdatedAt))

	var result domain.Wallet
	if err := query.QueryContext(ctx, r.db, &result); err != nil {
		return domain.Wallet{}, err
	}

	return result, nil
}

func (r *WalletRepo) List(ctx context.Context, account uuid.UUID) ([]domain.Wallet, error) {
	query := r.wallet.SELECT(
		r.wallet.AllColumns.Except(r.wallet.CreatedAt, r.wallet.UpdatedAt)).
		WHERE(r.wallet.Account.EQ(pg.UUID(account)))

	result := make([]domain.Wallet, 0, 1)
	if err := query.QueryContext(ctx, r.db, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *WalletRepo) Get(ctx context.Context, id int) (domain.Wallet, error) {
	query := r.wallet.SELECT(
		r.wallet.AllColumns.Except(r.wallet.CreatedAt, r.wallet.UpdatedAt)).
		WHERE(r.wallet.ID.EQ(pg.Int(int64(id))))

	var result domain.Wallet
	if err := query.QueryContext(ctx, r.db, &result); err != nil {
		return domain.Wallet{}, err
	}

	return result, nil
}

func (r *WalletRepo) HasTransaction(ctx context.Context, transaction domain.Transaction) (bool, error) {
	query := r.transaction.SELECT(pg.COUNT(r.transaction.WalletID).AS("count")).
		WHERE(r.transaction.WalletID.EQ(pg.Int(int64(transaction.WalletID))).
			AND(r.transaction.TransactionID.EQ(pg.UUID(transaction.ID))))

	var result struct {
		Count int
	}
	if err := query.QueryContext(ctx, r.db, &result); err != nil {
		return false, err
	}

	return result.Count > 0, nil
}

func (r *WalletRepo) ProcessTransaction(ctx context.Context, transaction domain.Transaction) (domain.Wallet, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return domain.Wallet{}, err
	}
	defer tx.Rollback()

	if err := r.createTransaction(ctx, tx, transaction); err != nil {
		return domain.Wallet{}, err
	}

	w, err := r.updateWallet(ctx, tx, transaction)
	if err != nil {
		return domain.Wallet{}, err
	}

	if err := tx.Commit(); err != nil {
		return domain.Wallet{}, err
	}
	return w, nil
}

func (r *WalletRepo) createTransaction(ctx context.Context, db qrm.Executable, transaction domain.Transaction) error {
	query := r.transaction.INSERT(r.transaction.WalletID, r.transaction.TransactionID).
		VALUES(transaction.WalletID, transaction.ID)

	if _, err := query.ExecContext(ctx, db); err != nil {
		return err
	}

	return nil
}

func (r *WalletRepo) updateWallet(ctx context.Context, db qrm.Queryable, transaction domain.Transaction) (domain.Wallet, error) {
	query := r.wallet.UPDATE(r.wallet.Amount).
		SET(r.wallet.Amount.ADD(pg.Int(int64(transaction.Amount)))).
		WHERE(r.wallet.ID.EQ(pg.Int(int64(transaction.WalletID))).
			AND(r.wallet.Currency.EQ(pg.String(string(transaction.Currency))))).
		RETURNING(r.wallet.AllColumns.Except(r.wallet.CreatedAt, r.wallet.UpdatedAt))

	var result domain.Wallet
	if err := query.QueryContext(ctx, db, &result); err != nil {
		return domain.Wallet{}, err
	}

	return result, nil
}
