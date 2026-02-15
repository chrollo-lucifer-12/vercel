package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/chrollo-lucifer-12/api-server/auth"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/utils"
)

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

	ctx := r.Context()

	_, err = h.db.GetUser(ctx, u.Email)
	if err == nil {
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	newUser := db.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: hashedPassword,
	}

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

	accessToken, accessClaims, err := h.auth.Maker.CreateToken(user.ID, user.Email, 15*time.Minute)

	if err != nil {
		http.Error(w, "Error creating token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken, refreshClaims, err := h.auth.Maker.CreateToken(user.ID, user.Email, 7*24*time.Hour)

	newSession := db.Session{
		UserID:       utils.StringToUUID(refreshClaims.RegisteredClaims.ID),
		UserEmail:    user.Email,
		RefreshToken: refreshToken,
		Revoked:      false,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	}
	h.auth.CreateSession(ctx, &newSession)

	res := LoginResponse{
		SessionID:             newSession.UserID.String(),
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		AccessToken:           accessToken,
		User: UserRes{
			Name:  user.Name,
			Email: user.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *ServerClient) logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	id := parts[2]
	sessionID := utils.StringToUUID(id)

	claims := r.Context().Value(authKey{}).(*auth.UserClaims)

	ctx := r.Context()

	session, err := h.auth.GetSession(ctx, sessionID)
	if err != nil {
		http.Error(w, "error getting session: "+err.Error(), http.StatusNotFound)
		return
	}

	if session.UserEmail != claims.Email {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	err = h.auth.DeleteSession(ctx, sessionID)
	if err != nil {
		http.Error(w, "failed to delete session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (h *ServerClient) refreshAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req RenewAccessTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.ErrBodyNotAllowed.Error(), http.StatusBadRequest)
		return
	}

	refreshClaims, err := h.auth.Maker.VerifyToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "error verifying refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	session, err := h.auth.GetSession(ctx, utils.StringToUUID(refreshClaims.RegisteredClaims.ID))
	if err != nil {
		http.Error(w, "error fetching session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if session.Revoked {
		http.Error(w, "session is revoked", http.StatusInternalServerError)
		return
	}

	if session.UserEmail != refreshClaims.Email {
		http.Error(w, "invalid session", http.StatusInternalServerError)
		return
	}

	accessToken, accessClaims, err := h.auth.Maker.CreateToken(refreshClaims.ID, refreshClaims.Email, 15*time.Minute)
	if err != nil {
		http.Error(w, "error creating access token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RenewAccessTokenRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	})
}
