package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/okoraretega/doc_stream_server/model"
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
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"user":   user,
		"wallet": wallet,
	})
}

func (wh *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()

	var wT model.WalletTransfer

	err := json.NewDecoder(r.Body).Decode(&wT)
	if err != nil {
		http.Error(w, "Unable to decode request", http.StatusBadRequest)
		return
	}

	wallet, err := wh.walletService.Transfer(ctx, wT.UserId, wT.FromAccount, wT.ToAccount, wT.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(wallet)

}
