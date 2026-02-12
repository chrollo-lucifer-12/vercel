package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/utils"
)

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRes struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	User        UserRes
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

func (h *ServerClient) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var l LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		http.Error(w, http.ErrBodyNotAllowed.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	user, err := h.auth.GetUser(ctx, l.Email)

	if err != nil {
		http.Error(w, "Error finding the user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ok := utils.CheckPassword(l.Password, user.Password)
	if !ok {
		http.Error(w, "Email or password is incorrect", http.StatusUnauthorized)
		return
	}

	accessToken, _, err := h.auth.Maker.CreateToken(user.ID, user.Email, 15*time.Minute)

	if err != nil {
		http.Error(w, "Error creating token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := LoginResponse{
		AccessToken: accessToken,
		User: UserRes{
			Name:  user.Name,
			Email: user.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
