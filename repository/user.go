package repository

import (
	"context"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/okoraretega/doc_stream_server/model"
)

type UserStore struct {
	users []model.User
	mu    sync.Mutex
}

type UserRepository interface {
	CreateUser(u model.User) (model.User, error)
	DeleteUser(id uuid.UUID) (bool, error)
	GetUserById(id uuid.UUID) (model.User, bool)
	GetAllUsers(ctx context.Context) ([]model.User, error)
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: []model.User{},
	}
}

func (s *UserStore) CreateUser(u model.User) (model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users = append(s.users, u)
	return u, nil
}

func (s *UserStore) GetAllUsers(ctx context.Context) ([]model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cpyUsers := make([]model.User, len(s.users))
	copy(cpyUsers, s.users)
	return cpyUsers, nil
}

func (s *UserStore) DeleteUser(id uuid.UUID) (bool, error) {
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

func (s *UserStore) GetUserById(id uuid.UUID) (model.User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, user := range s.users {
		if id == user.ID {
			return user, true
		}
	}

	return model.User{}, false
}
