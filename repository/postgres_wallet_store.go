package repository

import (
	"context"

	"github.com/okoraretega/doc_stream_server/model"
)

func (s *PostgresStore) GetAllWalets(ctx context.Context) ([]model.Wallet, error) {
	var wallets []model.Wallet
	rows, err := s.db.Query(ctx, "SELECT * FROM wallets")
	if err != nil {
		return []model.Wallet{}, err
	}

	defer rows.Close()
	for rows.Next() {
		var w model.Wallet
		err = rows.Scan(&w.Id, &w.WalletNumber, &w.WalletBalance, &w.UserId, &w.CreatedAt)
		if err != nil {
			return []model.Wallet{}, err
		}

		wallets = append(wallets, w)
	}
	return wallets, nil
}
