package services

import (
	"errors"
	"github.com/sirupsen/logrus"
	"url-shortener/internal/repositories"
	"url-shortener/pkg/models"
)

type UserService struct {
	Repo   *repositories.UserRepository
	Logger *logrus.Logger
}

func NewUserService(repos *repositories.UserRepository, logger *logrus.Logger) *UserService {
	return &UserService{Repo: repos, Logger: logger}
}

func (s *UserService) RegisterUser(user models.User) (int64, error) {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		s.Logger.Error("User registration failed, empty field(s)")
		return 0, errors.New("username or email or password is empty")
	}
	return s.Repo.InsertUser(user)
}

func (s *UserService) LoginUser(email string) (models.User, error) {
	if email == "" {
		s.Logger.Error("User login failed, empty field")
		return models.User{}, errors.New("email is empty")
	}
	return s.Repo.LoginUser(email)
}
