package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"transaction-system/internal/database"
	"transaction-system/internal/models"
)

type WalletHandler struct {
	db database.DatabaseInterface
}

func NewWalletHandler(db database.DatabaseInterface) *WalletHandler {
	return &WalletHandler{db}
}

func (h *WalletHandler) Wallets(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.db.Exec("UPDATE wallets SET balance = balance - $1 WHERE address = $2",
		req.Amount, req.From)
	if err != nil {
		http.Error(w, "Error updating sender wallet", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec("UPDATE wallets SET balance = balance + $1 WHERE address = $2",
		req.Amount, req.To)
	if err != nil {
		http.Error(w, "Error updating receiver wallet", http.StatusInternalServerError)
		return
	}

	if err = h.db.InsertTransaction(req.From, req.To, req.Amount); err != nil {
		http.Error(w, "Error logging transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Successfully transferred %.2f from %s to %s",
			req.Amount, req.From, req.To),
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
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 || pathParts[0] != "" ||
		pathParts[1] != "api" || pathParts[2] != "wallet" ||
		pathParts[4] != "balance" {
		errorResponse(w, "Invalid URL format. Use /api/wallet/{address}/balance",
			http.StatusBadRequest)
		return
	}

	address := pathParts[3]
	var balance float64
	err := h.db.QueryRow("SELECT balance FROM wallets WHERE address = $1",
		address).Scan(&balance)

	switch {
	case err == sql.ErrNoRows:
		errorResponse(w, "Wallet not found", http.StatusNotFound)
		return
	case err != nil:
		errorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := models.BalanceResponse{Balance: balance}
	json.NewEncoder(w).Encode(response)
}

func errorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error":  message,
		"status": strconv.Itoa(statusCode),
	})
}
