// cmd/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"transaction-system/internal/infrastructure/postgres"
	"transaction-system/internal/interfaces/handlers"
	"transaction-system/internal/interfaces/services"
	"transaction-system/pkg/utils"

	_ "github.com/lib/pq"
)

func main() {
	// Инициализация БД
	db, err := sql.Open("postgres", connString())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	// Инициализация репозиториев
	walletRepo := postgres.NewWalletRepository(db)
	transRepo := postgres.NewTransactionRepository(db)

	// Инициализация сервисов
	walletService := services.NewWalletService(walletRepo, transRepo)

	// Инициализация обработчиков
	handler := handlers.NewWalletHandler(walletService)

	// Создание таблиц
	if err := createTables(db); err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Инициализация кошельков
	if err := initializeWallets(db); err != nil {
		log.Fatal("Failed to initialize wallets:", err)
	}

	// Настройка маршрутов
	http.HandleFunc("/api/wallets", enableCORS(handler.Wallets))
	http.HandleFunc("/api/send", enableCORS(handler.Send))
	http.HandleFunc("/api/transactions", enableCORS(handler.GetLastTransactions))
	http.HandleFunc("/api/wallet/", enableCORS(handleWalletBalance(handler)))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func connString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func handleWalletBalance(h *handlers.WalletHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/balance") {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		h.GetBalance(w, r)
	}
}

func createTables(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS wallets (
			address VARCHAR PRIMARY KEY,
			balance FLOAT
		);
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			from_address VARCHAR NOT NULL,
			to_address VARCHAR NOT NULL,
			amount FLOAT NOT NULL,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`); err != nil {
		return err
	}
	return nil
}

func initializeWallets(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM wallets").Scan(&count); err != nil {
		return err
	}

	if count == 0 {
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		for i := 0; i < 10; i++ {
			address := utils.GenerateRandomAddress()
			if _, err := tx.Exec("INSERT INTO wallets (address, balance) VALUES ($1, $2)", address, 100.0); err != nil {
				tx.Rollback()
				return err
			}
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}
