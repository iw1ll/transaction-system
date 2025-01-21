package main

import (
	"log"
	"net/http"
	"transaction-system/internal/database"
	"transaction-system/internal/handlers"
)

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

	http.HandleFunc("/api/send", walletHandler.Send)
	http.HandleFunc("/api/transactions/", walletHandler.GetLastTransactions)
	http.HandleFunc("/api/wallet/", walletHandler.GetBalance)

	log.Println("Server is starting on port :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
