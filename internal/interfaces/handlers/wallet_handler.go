// internal/interfaces/handlers/wallet_handler.go
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"transaction-system/internal/domain"
	"transaction-system/internal/interfaces/services"
)

type WalletHandler struct {
	service *services.WalletService
}

func NewWalletHandler(service *services.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

// Wallets возвращает список всех кошельков
func (h *WalletHandler) Wallets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	wallets, err := h.service.GetAllWallets(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get wallets")
		return
	}

	respondJSON(w, http.StatusOK, wallets)
}

// Send обрабатывает перевод средств
func (h *WalletHandler) Send(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req domain.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	if err := h.service.TransferFunds(r.Context(), req); err != nil {
		switch {
		case errors.Is(err, domain.ErrWalletNotFound):
			respondError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrInsufficientFunds):
			respondError(w, http.StatusBadRequest, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "success",
	})
}

// GetLastTransactions возвращает последние транзакции
func (h *WalletHandler) GetLastTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	countStr := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		respondError(w, http.StatusBadRequest, "invalid count parameter")
		return
	}

	transactions, err := h.service.GetRecentTransactions(r.Context(), count)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get transactions")
		return
	}

	respondJSON(w, http.StatusOK, transactions)
}

// GetBalance возвращает баланс кошелька
func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 || pathParts[4] != "balance" {
		respondError(w, http.StatusBadRequest, "invalid URL format")
		return
	}

	address := pathParts[3]
	balance, err := h.service.GetWalletBalance(r.Context(), address)
	if err != nil {
		if errors.Is(err, domain.ErrWalletNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, domain.BalanceResponse{Balance: balance})
}

// Вспомогательные функции
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]interface{}{
		"error":  message,
		"status": status,
	})
}
