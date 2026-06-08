package services

import (
	"context"

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

func (s *WalletService) GetAllWalets(ctx context.Context) ([]model.Wallet, error) {
	return s.walletStore.GetAllWalets(ctx)
}
