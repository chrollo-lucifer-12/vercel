package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chrollo-lucifer-12/api-server/auth"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/workflow"
)

type Middleware func(http.Handler) http.Handler

type ServerClient struct {
	wClient *workflow.WorkflowClient
	db      *db.DB
	auth    *auth.AuthService
}

func NewServerClient(wClient *workflow.WorkflowClient, db *db.DB) (*ServerClient, error) {
	if wClient == nil {
		return nil, fmt.Errorf("No workflow client")
	}
	if db == nil {
		return nil, fmt.Errorf("No db client")
	}

	auth := auth.NewAuthService(auth.UserStoreFuncs{
		CreateUserFn: db.CreateUser,
		GetUserFn:    db.GetUser,
		UpdateUserFn: db.UpdateUser,
		DeleteUserFn: db.DeleteUser,
	}, auth.TokenStoreFuncs{
		CreateSessionFn: db.CreateSession,
		GetSessionFn:    db.GetSession,
		RevokeSessionFn: db.RevokeSession,
		DeleteSessionFn: db.DeleteSession,
	})

	return &ServerClient{wClient: wClient, db: db, auth: auth}, nil
}

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func (h *ServerClient) StartHTTP() {

	http.Handle(
		"/api/v1/deploy",
		Chain(
			http.HandlerFunc(h.deployHandler),
			h.authMiddleware,
		),
	)
	http.Handle("/api/v1/project",
		Chain(
			http.HandlerFunc(h.projectHandler),
			h.authMiddleware,
		),
	)
	http.Handle("/api/v1/auth/logout",
		Chain(
			http.HandlerFunc(h.logoutUserHandler),
			h.authMiddleware,
		),
	)
	http.HandleFunc("/api/v1/auth/register", h.registerUserHandler)
	http.HandleFunc("/api/v1/auth/login", h.loginUserHandler)
	http.HandleFunc("/api/v1/auth/refresh", h.refreshAccessTokenHandler)

	log.Fatal(http.ListenAndServe(":9000", nil))
}
