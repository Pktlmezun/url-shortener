package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"url-shortener/internal/services"
	"url-shortener/pkg/models"
)

type UserHandler struct {
	Service *services.UserService
	Logger  *logrus.Logger
}

func NewUserHandler(s *services.UserService, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		Service: s,
		Logger:  logger,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		h.Logger.Error("error while decoding request body")
		return
	}
	fmt.Println(user)
	hashedPassword, err := HashPassword(user.Password)
	fmt.Println(string(hashedPassword))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		h.Logger.Error("error while hashing password")
		return
	}
	h.Logger.Info("Successfully decoded user model")
	user.Password = string(hashedPassword)
	fmt.Println(user)
	user.Id, err = h.Service.RegisterUser(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&user)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		h.Logger.Error("error while decoding request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(user)
	fmt.Println(user.Password)
	hashedPassword, err := HashPassword(user.Password)
	fmt.Println(string(hashedPassword))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		h.Logger.Error("error while hashing password")
		return
	}
	h.Logger.Info("Successfully decoded user model")
	dbUser, err := h.Service.LoginUser(&user)
	if err != nil || dbUser.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.New("invalid credentials").Error()))
	}
	err = validatePassword(dbUser.Password, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&dbUser)

}

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func validatePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
