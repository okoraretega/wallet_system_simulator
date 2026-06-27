package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/okoraretega/doc_stream_server/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u model.User) (model.User, model.Wallet, error)
	DeleteUser(ctx context.Context, id uuid.UUID) (bool, error)
	GetUserById(ctx context.Context, id uuid.UUID) (model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, bool, error)
}

type WalletRepository interface {
	GetAllWallets(ctx context.Context) ([]model.Wallet, error)
	GetWalletByUserId(ctx context.Context, id uuid.UUID) (model.User, model.Wallet, error)
	Transfer(ctx context.Context, id uuid.UUID, fromAccount, toAccount string, amount float64) (model.Wallet, error)
}
