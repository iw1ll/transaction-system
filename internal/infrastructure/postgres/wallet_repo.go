package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"transaction-system/internal/domain"
)

// WalletRepo реализует интерфейс WalletRepository для PostgreSQL
type WalletRepo struct {
	db *sql.DB
}

// NewWalletRepository создает новый экземпляр WalletRepo
func NewWalletRepository(db *sql.DB) *WalletRepo {
	return &WalletRepo{db: db}
}

// GetAll возвращает все кошельки из базы данных
func (r *WalletRepo) GetAll(ctx context.Context) ([]domain.Wallet, error) {
	const query = `SELECT address, balance FROM wallets`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}
	defer rows.Close()

	var wallets []domain.Wallet
	for rows.Next() {
		var w domain.Wallet
		if err := rows.Scan(&w.Address, &w.Balance); err != nil {
			return nil, fmt.Errorf("failed to scan wallet: %w", err)
		}
		wallets = append(wallets, w)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return wallets, nil
}

// GetByAddress возвращает кошелек по адресу
func (r *WalletRepo) GetByAddress(ctx context.Context, address string) (*domain.Wallet, error) {
	const query = `SELECT address, balance FROM wallets WHERE address = $1`

	row := r.db.QueryRowContext(ctx, query, address)

	var wallet domain.Wallet
	if err := row.Scan(&wallet.Address, &wallet.Balance); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &wallet, nil
}

// Transfer выполняет перевод средств между кошельками в транзакции
func (r *WalletRepo) Transfer(ctx context.Context, from, to string, amount float64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Уменьшаем баланс отправителя
	updateFromQuery := `UPDATE wallets SET balance = balance - $1 WHERE address = $2`
	if _, err = tx.ExecContext(ctx, updateFromQuery, amount, from); err != nil {
		return fmt.Errorf("failed to update sender balance: %w", err)
	}

	// Увеличиваем баланс получателя
	updateToQuery := `UPDATE wallets SET balance = balance + $1 WHERE address = $2`
	if _, err = tx.ExecContext(ctx, updateToQuery, amount, to); err != nil {
		return fmt.Errorf("failed to update receiver balance: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Create создает новый кошелек
func (r *WalletRepo) Create(ctx context.Context, address string, balance float64) error {
	const query = `INSERT INTO wallets (address, balance) VALUES ($1, $2)`

	if _, err := r.db.ExecContext(ctx, query, address, balance); err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}

// Exists проверяет существование хотя бы одного кошелька
func (r *WalletRepo) Exists(ctx context.Context) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM wallets LIMIT 1)`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check wallets existence: %w", err)
	}

	return exists, nil
}
