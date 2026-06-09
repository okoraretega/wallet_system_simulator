package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/okoraretega/doc_stream_server/model"
)

func (s *PostgresStore) GetAllWallets(ctx context.Context) ([]model.Wallet, error) {
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

func (s *PostgresStore) GetWalletByUserId(ctx context.Context, id uuid.UUID) (model.User, model.Wallet, error) {
	var u model.User
	var wallet model.Wallet

	query := `SELECT users.first_name, users.last_name, wallets.id, wallets.wallet_number, wallets.wallet_balance, wallets.user_id, wallets.created_at FROM users JOIN wallets ON wallets.user_id = users.id WHERE users.id = $1`

	err := s.db.QueryRow(ctx, query, id).Scan(&u.FirstName, &u.LastName, &wallet.Id, &wallet.WalletNumber, &wallet.WalletBalance, &wallet.UserId, &wallet.CreatedAt)
	if err != nil {
		return model.User{}, model.Wallet{}, fmt.Errorf("Unable to scan user: %w", err)
	}
	return u, wallet, nil
}

func (s *PostgresStore) Transfer(ctx context.Context, user_id uuid.UUID, fromAccount, toAccount string, amount float64) (model.Wallet, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return model.Wallet{}, fmt.Errorf("Failed to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	debitQuery := `UPDATE wallets
	 SET wallet_balance = wallet_balance - $1
	 WHERE user_id = $2 AND wallet_balance >= $1`

	result, err := tx.Exec(ctx, debitQuery, amount, user_id)
	if err != nil {
		return model.Wallet{}, fmt.Errorf("Failed to debit user: %w", err)
	}

	if result.RowsAffected() < 1 {
		return model.Wallet{}, fmt.Errorf("Unable to debit user wallet")
	}

	creditQuery := `UPDATE wallets
	SET wallet_balance = wallet_balance + $1
	WHERE wallet_number = $2`

	result, err = tx.Exec(ctx, creditQuery, amount, toAccount)
	if err != nil {
		return model.Wallet{}, fmt.Errorf("Failed to credit user: %w", err)
	}

	if result.RowsAffected() < 1 {
		return model.Wallet{}, fmt.Errorf("credit failed for wallet_number=%s: no rows affected", toAccount)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.Wallet{}, fmt.Errorf("Failed to commit transaction: %w", err)
	}

	_, wallet, err := s.GetWalletByUserId(ctx, user_id)
	if err != nil {
		return model.Wallet{}, fmt.Errorf("Failed to get user wallet: %w", err)
	}

	return wallet, nil
}
