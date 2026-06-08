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
