package database

import (
	"database/sql"
	"transaction-system/internal/models"
)

type DatabaseInterface interface {
	CreateTable() error
	CreateTransactionTable() error
	WalletsExist() bool
	CreateWallets(count int) error
	GetLastTransactions(count int) ([]models.Transaction, error)
	InsertTransaction(from, to string, amount float64) error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}
