package repository

import (
	"context"
	"errors"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/okoraretega/doc_stream_server/model"
)

type UserStore struct {
	users []model.User
	mu    sync.Mutex
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: []model.User{},
	}
}

func (s *UserStore) CreateUser(ctx context.Context, u model.User) (model.User, model.Wallet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users = append(s.users, u)
	return u, model.Wallet{}, nil
}

func (s *UserStore) GetAllUsers(ctx context.Context) ([]model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cpyUsers := make([]model.User, len(s.users))
	copy(cpyUsers, s.users)
	return cpyUsers, nil
}

func (s *UserStore) DeleteUser(ctx context.Context, id uuid.UUID) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, user := range s.users {
		if id == user.ID {
			s.users = slices.Delete(s.users, i, i+1)
			return true, nil
		}
	}

	return false, nil
}

func (s *UserStore) GetUserById(ctx context.Context, id uuid.UUID) (model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, user := range s.users {
		if id == user.ID {
			return user, nil
		}
	}

	return model.User{}, errors.New("Unable to get user")
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (model.User, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, user := range s.users {
		if user.Email == email {
			return user, true, nil
		}
	}
	return model.User{}, false, errors.New("Failed to find user")
}
