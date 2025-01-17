package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"transaction-system/internal/models"
)

type WalletHandler struct {
	db *sql.DB
}

func NewWalletHandler(db *sql.DB) *WalletHandler {
	return &WalletHandler{db}
}

func (h *WalletHandler) Send(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req models.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var balance float64
	err := h.db.QueryRow("SELECT balance FROM wallets WHERE address = $1", req.From).Scan(&balance)
	if err != nil {
		http.Error(w, "Sender wallet not found", http.StatusNotFound)
		return
	}

	if balance < req.Amount {
		http.Error(w, "Insufficient funds", http.StatusPaymentRequired)
		return
	}

	if _, err := h.db.Exec("UPDATE wallets SET balance = balance - $1 WHERE address = $2", req.Amount, req.From); err != nil {
		http.Error(w, "Error updating sender wallet", http.StatusInternalServerError)
		return
	}

	if _, err := h.db.Exec("UPDATE wallets SET balance = balance + $1 WHERE address = $2", req.Amount, req.To); err != nil {
		http.Error(w, "Error updating receiver wallet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Successfully transferred %.2f from %s to %s", req.Amount, req.From, req.To),
	})
}

func (h *WalletHandler) GetWallets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {

		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rows, err := h.db.Query("SELECT address, balance FROM wallets")
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
