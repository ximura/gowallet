package domain

import "github.com/google/uuid"

type Currency string

type Transaction struct {
	ID       uuid.UUID
	WalletID int
	Amount   int
	Currency Currency
}

type Wallet struct {
	ID       int
	Account  uuid.UUID
	Amount   int
	Currency Currency
}
