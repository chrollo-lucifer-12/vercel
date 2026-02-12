package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/utils"
)

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *ServerClient) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var u UserRequest
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, http.ErrBodyNotAllowed.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		http.Error(w, "Error hashing password: "+err.Error(), http.StatusInternalServerError)
	}

	newUser := db.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: hashedPassword,
	}

	ctx := context.Background()

	err = h.auth.CreateUser(ctx, &newUser)
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
