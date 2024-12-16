package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"url-shortener/internal/services"
	"url-shortener/pkg/models"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{
		Service: s,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	fmt.Println(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.Id, err = h.Service.RegisterUser(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&user)

}

//type validateUser struct {
//	Email    *string `json:"email"`
//	Password *string `json:"password"`
//	Username *string `json:"username"`
//}
