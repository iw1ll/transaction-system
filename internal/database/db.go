package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUsername,
		dbPassword,
		dbDatabase,
		dbHost,
		dbPort,
	)

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

	return tx.Commit()
}
