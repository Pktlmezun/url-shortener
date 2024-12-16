package services

import (
	"errors"
	"url-shortener/internal/repositories"
	"url-shortener/pkg/models"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func NewUserService(repos *repositories.UserRepository) *UserService {
	return &UserService{Repo: repos}
}

func (s *UserService) RegisterUser(user models.User) (int64, error) {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return 0, errors.New("username or email or password is empty")
	}

	return s.Repo.InsertUser(user)
}
