package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"transaction-system/internal/database"
	"transaction-system/internal/models"
)

type WalletHandler struct {
	db *database.Database
}

func NewWalletHandler(db *database.Database) *WalletHandler {
	return &WalletHandler{db}
}

func (h *WalletHandler) Send(w http.ResponseWriter, r *http.Request) {
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
	err = h.db.QueryRow("SELECT balance FROM wallets WHERE address = $1", req.From).Scan(&balance)
	if err != nil {
		http.Error(w, "Sender wallet not found", http.StatusNotFound)
		return
	}

	if balance < req.Amount {
		http.Error(w, "Insufficient funds", http.StatusForbidden)
		return
	}

	_, err = h.db.Exec("UPDATE wallets SET balance = balance - $1 WHERE address = $2", req.Amount, req.From)
	if err != nil {
		http.Error(w, "Error updating sender wallet", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec("UPDATE wallets SET balance = balance + $1 WHERE address = $2", req.Amount, req.To)
	if err != nil {
		http.Error(w, "Error updating receiver wallet", http.StatusInternalServerError)
		return
	}

	if err = h.db.InsertTransaction(req.From, req.To, req.Amount); err != nil {
		http.Error(w, "Error logging transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Successfully transferred %.2f from %s to %s", req.Amount, req.From, req.To),
	})
}

func (h *WalletHandler) GetLastTransactions(w http.ResponseWriter, r *http.Request) {
	countStr := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		http.Error(w, "Invalid count parameter", http.StatusBadRequest)
		return
	}

	transactions, err := h.db.GetLastTransactions(count)
	if err != nil {
		http.Error(w, "Error retrieving transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Path[len("/api/wallet/") : len(r.URL.Path)-len("/balance")]

	var balance float64
	err := h.db.QueryRow("SELECT balance FROM wallets WHERE address = $1", address).Scan(&balance)
	if err == sql.ErrNoRows {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error retrieving balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := models.BalanceResponse{Balance: balance}
	json.NewEncoder(w).Encode(response)
}
