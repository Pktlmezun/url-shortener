package handlers

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"url-shortener/internal/auth"
	"url-shortener/internal/services"
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
	user, err := auth.DecodeUser(r, h.Logger)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		h.Logger.Error("error while hashing password")
		return
	}
	user.Password = string(hashedPassword)
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
	user, err := auth.DecodeUser(r, h.Logger)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	dbUser, err := h.Service.LoginUser(user.Email)
	if err != nil || dbUser.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.New("invalid credentials").Error()))
	}
	err = auth.ValidatePassword(dbUser.Password, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := auth.CreateToken(dbUser.Id, dbUser.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logger.Errorf("Failed to create token: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	//_ = json.NewEncoder(w).Encode(&dbUser)
	response := map[string]string{
		"token": token,
	}
	_ = json.NewEncoder(w).Encode(response)
}
