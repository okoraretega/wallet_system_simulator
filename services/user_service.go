package services

import (
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

func (s *UserService) CreateUser(u model.User) (error, model.User) {
	users := s.userStore.GetAllUsers()
	for _, user := range users {
		if u.Email == user.Email {
			return errors.New("User with email already exists"), model.User{}
		}
	}

	s.userStore.CreateUser(u)
	return nil, u

}

func (s *UserService) GetUserById(id uuid.UUID) (model.User, bool) {
	return s.userStore.GetUserById(id)
}

func (s *UserService) DeleteUser(id uuid.UUID) bool {
	return s.userStore.DeleteUser(id)
}

func (s *UserService) GetAllUsers() []model.User {
	return s.userStore.GetAllUsers()
}
