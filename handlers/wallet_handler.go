package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/okoraretega/doc_stream_server/services"
)

type WalletHandler struct {
	walletService *services.WalletService
}

func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

func (wh *WalletHandler) GetAllWalets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wallets, err := wh.walletService.GetAllWalets(ctx)
	if err != nil {
		http.Error(w, "Failed to get wallets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(wallets)
}
