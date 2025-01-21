package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"transaction-system/internal/models"
	"transaction-system/internal/utils"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

func NewDatabase() (*Database, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUsername := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbDatabase := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUsername, dbPassword, dbDatabase, dbHost, dbPort)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (d *Database) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS wallets (
		address VARCHAR PRIMARY KEY,
		balance FLOAT
	);`
	_, err := d.Exec(query)
	return err
}

func (d *Database) CreateTransactionTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			from_address VARCHAR NOT NULL,
			to_address VARCHAR NOT NULL,
			amount FLOAT NOT NULL,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`
	_, err := d.Exec(query)
	return err
}

func (d *Database) WalletsExist() bool {
	var count int
	row := d.QueryRow("SELECT COUNT(*) FROM wallets")
	err := row.Scan(&count)
	return err == nil && count > 0
}

func (d *Database) CreateWallets(count int) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		address := utils.GenerateRandomAddress()
		if _, err := tx.Exec("INSERT INTO wallets (address, balance) VALUES ($1, $2)", address, 100.0); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *Database) GetLastTransactions(count int) ([]models.Transaction, error) {
	rows, err := d.Query("SELECT from_address, to_address, amount, timestamp FROM transactions ORDER BY timestamp DESC LIMIT $1", count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.From, &t.To, &t.Amount, &t.Timestamp); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (d *Database) InsertTransaction(from, to string, amount float64) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}

	if _, err := d.Exec("INSERT INTO transactions (from_address, to_address, amount) VALUES ($1, $2, $3)", from, to, amount); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
