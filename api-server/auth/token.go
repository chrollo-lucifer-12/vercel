package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey: secretKey}
}

type UserClaims struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	jwt.RegisteredClaims
}

func NewUserClaims(id uuid.UUID, email string, duration time.Duration) (*UserClaims, error) {
	tokenId, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	return &UserClaims{
		ID:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}

func (j *JWTMaker) CreateToken(id uuid.UUID, email string, duration time.Duration) (string, *UserClaims, error) {
	claims, err := NewUserClaims(id, email, duration)
	if err != nil {
		return "", nil, fmt.Errorf("error creating claims: %w", err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(j.secretKey)

	if err != nil {
		return "", nil, fmt.Errorf("error signing tokne: %w", err)
	}

	return tokenStr, claims, nil
}

func (j *JWTMaker) VerifyToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Invalid token signing method")
		}

		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("Invalid claims")
	}

	return claims, nil
}
