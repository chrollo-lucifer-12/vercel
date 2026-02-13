package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/chrollo-lucifer-12/api-server/auth"
)

func verifyClaimsFromHeader(r *http.Request, maker *auth.JWTMaker) (*auth.UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("Authorization header is missing")
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 2 {
		return nil, fmt.Errorf("Invalid authorization header")
	}

	token := fields[1]
	return maker.VerifyToken(token)
}

func (h *ServerClient) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := verifyClaimsFromHeader(r, h.auth.Maker)
		if err != nil {
			http.Error(w, fmt.Errorf("error verifying token: %v", err).Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), authKey{}, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *ServerClient) methodMiddleware(method string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
