package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/okoraretega/doc_stream_server/model"
	"github.com/okoraretega/doc_stream_server/repository"
)

type WalletService struct {
	walletStore repository.WalletRepository
}

func NewWalletService(ws repository.WalletRepository) *WalletService {
	return &WalletService{
		walletStore: ws,
	}
}

func (s *WalletService) GetAllWallets(ctx context.Context) ([]model.Wallet, error) {
	return s.walletStore.GetAllWallets(ctx)
}

func (s *WalletService) GetWalletByUserId(ctx context.Context, id uuid.UUID) (model.User, model.Wallet, error) {
	return s.walletStore.GetWalletByUserId(ctx, id)
}

func (s *WalletService) Transfer(ctx context.Context, user_id uuid.UUID, fromAccount, toAccount string, amount float64) (model.Wallet, error) {
	if user_id == uuid.Nil && fromAccount == "" && toAccount == "" && amount > 1.00 {
		return model.Wallet{}, errors.New("User ID, Amount, Source account, Destination account and Amount are required")
	}
	return s.walletStore.Transfer(ctx, user_id, fromAccount, toAccount, amount)
}
