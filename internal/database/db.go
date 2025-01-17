package database

import (
	"database/sql"
	"log"
	"transaction-system/internal/utils"

	_ "github.com/lib/pq"
)

func InitDB(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to the database:", err)
	}

	return db
}

func CreateTable(db *sql.DB) {
	query := `
        CREATE TABLE IF NOT EXISTS wallets (
            address VARCHAR PRIMARY KEY,
            balance FLOAT
        );`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
}

func WalletsExist(db *sql.DB) bool {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM wallets")
	err := row.Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func CreateWallets(db *sql.DB, count int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < count; i++ {
		address := utils.GenerateRandomAddress()
		_, err := tx.Exec("INSERT INTO wallets (address, balance) VALUES ($1, $2)", address, 100.0)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created wallets:", count)
}
