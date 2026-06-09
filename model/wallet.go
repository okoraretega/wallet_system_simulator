package model

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	Id            uuid.UUID `json:"id"`
	WalletNumber  string    `json:"wallet_number"`
	WalletBalance float64   `json:"wallet_balance"`
	UserId        uuid.UUID `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type WalletTransfer struct {
	UserId      uuid.UUID `json:"user_id"`
	FromAccount string    `json:"from_accout"`
	ToAccount   string    `json:"to_account"`
	Amount      float64   `json:"amount"`
}
