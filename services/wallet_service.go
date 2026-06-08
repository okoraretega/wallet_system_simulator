package services

import (
	"context"

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
