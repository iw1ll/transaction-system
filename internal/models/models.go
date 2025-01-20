package models

import "time"

type Wallet struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}

type TransferRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type Transaction struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}
