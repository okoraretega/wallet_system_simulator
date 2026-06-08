package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
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
	wallets, err := wh.walletService.GetAllWallets(ctx)
	if err != nil {
		http.Error(w, "Failed to get wallets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(wallets)
}

func (wh *WalletHandler) GetWalletByUserId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	url := strings.TrimPrefix(r.URL.Path, "/wallets/")
	id, err := uuid.Parse(url)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusBadRequest)
		return
	}

	user, wallet, err := wh.walletService.GetWalletByUserId(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(map[string]any{
		"user":   user,
		"wallet": wallet,
	})
}
