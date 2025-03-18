package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"transaction-system/internal/domain"
	_ "transaction-system/internal/domain"
)

type TransactionRepo struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) Create(ctx context.Context, from, to string, amount float64) error {
	query := `INSERT INTO transactions (from_address, to_address, amount) 
              VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, from, to, amount)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepo) GetRecent(ctx context.Context, limit int) ([]domain.Transaction, error) {
	query := `
        SELECT from_address, to_address, amount, timestamp 
        FROM transactions 
        ORDER BY timestamp DESC 
        LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent transactions: %w", err)
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(&t.From, &t.To, &t.Amount, &t.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return transactions, nil
}
