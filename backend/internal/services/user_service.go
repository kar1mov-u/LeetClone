package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kar1mov-u/LeetClone/internal/models"
	"github.com/kar1mov-u/LeetClone/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repo.UserRepository
}

func NewUserService(userRepo *repo.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

var (
	EmailTakenErr    = errors.New("Email is already in use")
	UsernameTakenErr = errors.New("Username is already in use")
)

func (s *UserService) RegisterUser(context context.Context, data models.UserRegister) (uuid.UUID, error) {
	//check if the username or
	if !s.userRepo.CheckEmail(context, data.Email) {
		return uuid.UUID{}, EmailTakenErr
	}
	if !s.userRepo.CheckUsername(context, data.Username) {
		return uuid.UUID{}, UsernameTakenErr
	}
	data.Password, _ = hashPass(data.Password)
	userID, err := s.userRepo.CreateUser(context, data)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Fail to regster user: %v", err)
	}
	return userID, nil

}

func hashPass(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func verifyPass(hash, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil

}
