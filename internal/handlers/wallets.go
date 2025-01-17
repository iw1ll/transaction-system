package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"transaction-system/internal/models"
)

func WalletsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Неверный метод запроса", http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.Query("SELECT address, balance FROM wallets")
		if err != nil {
			http.Error(w, "Error fetching wallets", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var wallets []models.Wallet
		for rows.Next() {
			var wallet models.Wallet
			if err := rows.Scan(&wallet.Address, &wallet.Balance); err != nil {
				http.Error(w, "Error scanning wallet", http.StatusInternalServerError)
				return
			}
			wallets = append(wallets, wallet)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wallets)
	}
}
