package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/okoraretega/doc_stream_server/model"
	"github.com/okoraretega/doc_stream_server/repository"
	"golang.org/x/crypto/bcrypt"
)

var JWT_SECRET = "edwed32eedqEdDEDrwdweD" //This should be in your env

type UserService struct {
	userStore repository.UserRepository
}

func NewUserService(s repository.UserRepository) *UserService {
	return &UserService{
		userStore: s,
	}
}

func (s *UserService) CreateUser(ctx context.Context, u model.CreateUser) (model.User, model.Wallet, error) {

	user, ok, _ := s.GetUserByEmail(ctx, u.Email)

	if ok {
		if u.Email == user.Email {
			return model.User{}, model.Wallet{}, errors.New("A user with that email already exists")
		}

	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return model.User{}, model.Wallet{}, errors.New("Failed to encrypt password")
	}

	createdUser := model.User{
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email,
		PasswordHash: string(passwordHash),
	}

	newUser, wallet, err := s.userStore.CreateUser(ctx, createdUser)
	if err != nil {
		return model.User{}, model.Wallet{}, err
	}
	return newUser, wallet, nil

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

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (model.User, bool, error) {
	return s.userStore.GetUserByEmail(ctx, email)
}

func (s *UserService) Login(ctx context.Context, email, password string) (model.UserResponse, string, error) {

	user, ok, err := s.GetUserByEmail(ctx, email)
	if !ok && err != nil {
		return model.UserResponse{}, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return model.UserResponse{}, "", fmt.Errorf("Failed to compare hash: %w", err)
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	userRes := model.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(JWT_SECRET))

	return userRes, tokenStr, nil
}
