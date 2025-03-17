package handlers

import "net/http"

type WalletHandlerInterface interface {
	Wallets(w http.ResponseWriter, r *http.Request)
	Send(w http.ResponseWriter, r *http.Request)
	GetLastTransactions(w http.ResponseWriter, r *http.Request)
	GetBalance(w http.ResponseWriter, r *http.Request)
}
