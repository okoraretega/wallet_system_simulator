package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/okoraretega/doc_stream_server/helpers"
	"github.com/okoraretega/doc_stream_server/model"
)

func (s *PostgresStore) CreateUser(ctx context.Context, u model.User) (model.User, model.Wallet, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return model.User{}, model.Wallet{}, err
	}

	defer tx.Rollback(ctx)

	userQuery := `INSERT INTO users (first_name, last_name, email)
				VALUES($1, $2, $3)
				RETURNING id, created_at
	`

	err = tx.QueryRow(ctx, userQuery, u.FirstName, u.LastName, u.Email).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		return model.User{}, model.Wallet{}, err
	}

	walletNumber := helpers.GenerateWalletNumber()
	var w model.Wallet

	walletQuery := `INSERT INTO wallets (user_id, wallet_number)
					VALUES($1, $2)
					RETURNING id, wallet_number, user_id
	`

	err = tx.QueryRow(ctx, walletQuery, u.ID, walletNumber).Scan(&w.Id, &w.WalletNumber, &w.UserId)
	if err != nil {
		return model.User{}, model.Wallet{}, fmt.Errorf("Falied to create wallet for user: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.User{}, model.Wallet{}, fmt.Errorf("Failed to compelete user creation: %w", err)
	}
	return u, w, nil
}

func (s *PostgresStore) GetAllUsers(ctx context.Context) ([]model.User, error) {

	query := `SELECT id, first_name, last_name, email, created_at FROM users`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return []model.User{}, fmt.Errorf("Unable to query rows: %w", err)
	}

	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt)
		if err != nil {
			return []model.User{}, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (s *PostgresStore) GetUserById(ctx context.Context, id uuid.UUID) (model.User, error) {
	query := `SELECT id, first_name, last_name, email, created_at FROM users WHERE id = $1`

	var u model.User
	err := s.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("User not found: %w", err)
	}
	return u, err
}

func (s *PostgresStore) DeleteUser(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `DELETE FROM users WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}
