package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/okoraretega/doc_stream_server/model"
)

type PostgresUserStore struct {
	db *pgxpool.Pool
}

func NewPostgresUserStore(connString string) (*PostgresUserStore, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse db config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Minute * 30

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}

	return &PostgresUserStore{db: db}, nil
}

func (s *PostgresUserStore) Close() {
	s.db.Close()
}

func (s *PostgresUserStore) CreateUser(ctx context.Context, u model.User) (model.User, error) {

	query := `INSERT INTO users (first_name, last_name, email)
				VALUES($1, $2, $3)
				RETURNING id
	`
	err := s.db.QueryRow(ctx, query, u.FirstName, u.LastName, u.Email).Scan(&u.ID)
	if err != nil {
		return model.User{}, err
	}

	return u, nil
}

func (s *PostgresUserStore) GetAllUsers(ctx context.Context) ([]model.User, error) {

	query := `SELECT id, first_name, last_name, email FROM users`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		fmt.Printf("Unable to query rows")
	}

	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email)
		if err != nil {
			return []model.User{}, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (s *PostgresUserStore) GetUserById(ctx context.Context, id uuid.UUID) (model.User, bool) {
	query := `SELECT id, first_name, last_name, email FROM users WHERE id = $1`

	var u model.User
	err := s.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email)
	if err != nil {
		fmt.Printf("Unable to get user from db")
		return model.User{}, false
	}
	return u, true
}

func (s *PostgresUserStore) DeleteUser(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `DELETE FROM users WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}
