package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/okoraretega/doc_stream_server/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u model.User) (model.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) (bool, error)
	GetUserById(ctx context.Context, id uuid.UUID) (model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
}

type WalletRepository interface {
	GetAllWalets(ctx context.Context) ([]model.Wallet, error)
}
