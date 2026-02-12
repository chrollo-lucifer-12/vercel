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

type TokenStoreFuncs struct {
	CreateSessionFn func(ctx context.Context, s *db.Session) error
	GetSessionFn    func(ctx context.Context, id uuid.UUID) (*db.Session, error)
	RevokeSessionFn func(ctx context.Context, s db.Session) error
	DeleteSessionFn func(ctx context.Context, id uuid.UUID) error
}

type AuthService struct {
	userStore  UserStoreFuncs
	Maker      *JWTMaker
	tokenStore TokenStoreFuncs
}

func NewAuthService(userStore UserStoreFuncs, tokenStore TokenStoreFuncs) *AuthService {

	maker := NewJWTMaker("jowfuuf")

	return &AuthService{
		userStore:  userStore,
		Maker:      maker,
		tokenStore: tokenStore,
	}
}

func (a *AuthService) CreateSession(ctx context.Context, s *db.Session) error {
	return a.tokenStore.CreateSessionFn(ctx, s)
}

func (a *AuthService) GetSession(ctx context.Context, id uuid.UUID) (*db.Session, error) {
	return a.tokenStore.GetSessionFn(ctx, id)
}

func (a *AuthService) RevokeSession(ctx context.Context, s db.Session) error {
	return a.tokenStore.RevokeSessionFn(ctx, s)
}

func (a *AuthService) DeleteSession(ctx context.Context, id uuid.UUID) error {
	return a.tokenStore.DeleteSessionFn(ctx, id)
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
