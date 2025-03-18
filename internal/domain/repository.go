// internal/domain/repository.go
package domain

import "context"

type WalletRepository interface {
	GetAll(ctx context.Context) ([]Wallet, error)
	GetByAddress(ctx context.Context, address string) (*Wallet, error)
	Transfer(ctx context.Context, from, to string, amount float64) error
	Create(ctx context.Context, address string, balance float64) error
	Exists(ctx context.Context) (bool, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, from, to string, amount float64) error
	GetRecent(ctx context.Context, limit int) ([]Transaction, error)
}
