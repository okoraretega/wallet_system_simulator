package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/okoraretega/doc_stream_server/model"
	"github.com/okoraretega/doc_stream_server/repository"
)

type UserService struct {
	userStore repository.UserRepository
}

func NewUserService(s repository.UserRepository) *UserService {
	return &UserService{
		userStore: s,
	}
}

func (s *UserService) CreateUser(ctx context.Context, u model.User) (model.User, error) {
	users, err := s.userStore.GetAllUsers(ctx)
	if err != nil {
		return model.User{}, errors.New("An error occured while reading all users")
	}
	for _, user := range users {
		if u.Email == user.Email {
			return model.User{}, errors.New("User with email already exists")
		}
	}

	newUser, err := s.userStore.CreateUser(ctx, u)
	if err != nil {
		return model.User{}, err
	}
	return newUser, nil

}

func (s *UserService) GetUserById(ctx context.Context, id uuid.UUID) (model.User, error) {
	return s.userStore.GetUserById(ctx, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.userStore.DeleteUser(ctx, id)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return s.userStore.GetAllUsers(ctx)
}
