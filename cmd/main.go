package main

import (
	"log"
	"net/http"
	"transaction-system/internal/database/database"
	"transaction-system/internal/database/handlers"
)

func main() {
	db := database.InitDB("user=user password=password dbname=transaction_system host=localhost port=5432 sslmode=disable")
	defer db.Close()

	database.CreateTable(db)

	if !database.WalletsExist(db) {
		database.CreateWallets(db, 10)
	}

	http.HandleFunc("/api/send", handlers.SendHandler(db))
	http.HandleFunc("/api/wallets", handlers.WalletsHandler(db))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
