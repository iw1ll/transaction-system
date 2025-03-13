package main

import (
	"log"
	"net/http"
	"strings"
	"transaction-system/internal/database"
	"transaction-system/internal/handlers"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Successfully connected to the database")

	if err := db.CreateTable(); err != nil {
		log.Fatal(err)
	}

	if err := db.CreateTransactionTable(); err != nil {
		log.Fatal(err)
	}

	if !db.WalletsExist() {
		if err := db.CreateWallets(10); err != nil {
			log.Fatal(err)
		}
	}

	walletHandler := handlers.NewWalletHandler(db)

	http.HandleFunc("/api/wallets", enableCORS(walletHandler.Wallets))
	http.HandleFunc("/api/send", enableCORS(walletHandler.Send))
	http.HandleFunc("/api/transactions", enableCORS(walletHandler.GetLastTransactions))
	http.HandleFunc("/api/wallet/", enableCORS(handleWalletBalance(walletHandler)))

	log.Println("Server is starting on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
