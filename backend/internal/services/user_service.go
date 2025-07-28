package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kar1mov-u/LeetClone/internal/models"
	"github.com/kar1mov-u/LeetClone/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo  *repo.UserRepository
	JwtSecret string
}

func NewUserService(userRepo *repo.UserRepository, jwtSecret string) *UserService {
	return &UserService{userRepo: userRepo, JwtSecret: jwtSecret}
}

var (
	EmailTakenErr         = errors.New("Email is already in use")
	UsernameTakenErr      = errors.New("Username is already in use")
	InvalidCredentialsErr = errors.New("Invalid credentials")
	UserNotFoundErr       = errors.New("User not found")
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

func (s *UserService) LoginUser(context context.Context, data models.UserLogin) (string, error) {
	//get password from the DB
	userID, dbPass, err := s.userRepo.GetUserPassword(context, data.Username)
	if err != nil {
		log.Println(err)

		return "", UserNotFoundErr
	}

	if !verifyPass(dbPass, data.Password) {
		return "", InvalidCredentialsErr
	}

	//create jwt
	accessToken, err := createToken(userID, 30, s.JwtSecret)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *UserService) GetUserByID(context context.Context, id uuid.UUID) (models.User, error) {
	return s.userRepo.GetUserByID(context, id)

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

func createToken(userID uuid.UUID, expiresMinutes int, key string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(time.Duration(expiresMinutes) * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
