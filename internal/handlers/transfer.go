package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"transaction-system/internal/models"
)

func SendHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {

			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var req models.TransferRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		var balance float64
		err = db.QueryRow("SELECT balance FROM wallets WHERE address = $1", req.From).Scan(&balance)
		if err != nil {
			http.Error(w, "Sender wallet not found", http.StatusNotFound)
			return
		}

		if balance < req.Amount {
			http.Error(w, "Insufficient funds", http.StatusPaymentRequired)
			return
		}

		_, err = db.Exec("UPDATE wallets SET balance = balance - $1 WHERE address = $2", req.Amount, req.From)
		if err != nil {
			http.Error(w, "Error updating sender wallet", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("UPDATE wallets SET balance = balance + $1 WHERE address = $2", req.Amount, req.To)
		if err != nil {
			http.Error(w, "Error updating receiver wallet", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("Successfully transferred %.2f from %s to %s", req.Amount, req.From, req.To),
		})
	}
}
