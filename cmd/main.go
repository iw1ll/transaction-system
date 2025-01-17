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

	if err := db.CreateTable(); err != nil {
		log.Fatal(err)
	}

	if !db.WalletsExist() {
		if err := db.CreateWallets(10); err != nil {
			log.Fatal(err)
		}
	}

	walletHandler := handlers.NewWalletHandler(db.DB)

	http.HandleFunc("/api/send", walletHandler.Send)
	http.HandleFunc("/api/wallets", walletHandler.GetWallets)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
