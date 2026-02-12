package auth

import (
	"context"

	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/google/uuid"
)

type UserStoreFuncs struct {
	CreateUserFn func(ctx context.Context, u *db.User) error
	GetUserFn    func(ctx context.Context, email string) (db.User, error)
	UpdateUserFn func(ctx context.Context, u db.User) error
	DeleteUserFn func(ctx context.Context, id uuid.UUID) error
}

type AuthService struct {
	userStore UserStoreFuncs
	Maker     *JWTMaker
}

func NewAuthService(userStore UserStoreFuncs) *AuthService {

	maker := NewJWTMaker("jowfuuf")

	return &AuthService{
		userStore: userStore,
		Maker:     maker,
	}
}

func (a *AuthService) CreateUser(ctx context.Context, u *db.User) error {
	return a.userStore.CreateUserFn(ctx, u)
}

func (a *AuthService) GetUser(ctx context.Context, email string) (db.User, error) {
	return a.userStore.GetUserFn(ctx, email)
}

func (a *AuthService) UpdateUser(ctx context.Context, u db.User) error {
	return a.userStore.UpdateUserFn(ctx, u)
}

func (a *AuthService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return a.userStore.DeleteUserFn(ctx, id)
}
