package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chrollo-lucifer-12/api-server/auth"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/mail"
	"github.com/chrollo-lucifer-12/shared/workflow"
)

func NewServerClient(wClient *workflow.WorkflowClient, dbClient *db.DB, mailClient *mail.MailClient) (*ServerClient, error) {

	if wClient == nil {
		return nil, fmt.Errorf("workflow client required")
	}

	if dbClient == nil {
		return nil, fmt.Errorf("db client required")
	}

	authService := newAuthService(dbClient)

	server := &ServerClient{
		wClient: wClient,
		db:      dbClient,
		auth:    authService,
		mail:    mailClient,
	}

	server.setupHTTP()

	return server, nil
}

func newAuthService(dbClient *db.DB) *auth.AuthService {

	return auth.NewAuthService(
		auth.UserStoreFuncs{
			CreateUserFn: dbClient.CreateUser,
			GetUserFn:    dbClient.GetUser,
			UpdateUserFn: dbClient.UpdateUser,
			DeleteUserFn: dbClient.DeleteUser,
		},
		auth.TokenStoreFuncs{
			CreateSessionFn: dbClient.CreateSession,
			GetSessionFn:    dbClient.GetSession,
			RevokeSessionFn: dbClient.RevokeSession,
			DeleteSessionFn: dbClient.DeleteSession,
		},
	)
}

func (s *ServerClient) setupHTTP() {
	logger := log.New(os.Stdout, "HTTP: ", log.LstdFlags)
	mux := http.NewServeMux()
	s.registerRoutes(mux)

	s.server = &http.Server{
		Addr:     ":9000",
		Handler:  enableCORS(mux),
		ErrorLog: logger,
	}
}

func (s *ServerClient) registerRoutes(mux *http.ServeMux) {

	routes := []route{

		{"/api/v1/deploy/create", http.MethodPost, s.deployHandler, true},
		{"/api/v1/project/create", http.MethodPost, s.createProjectHandler, true},
		{"/api/v1/projects", http.MethodGet, s.getAllProjectsHandler, true},
		{"/api/v1/project/", http.MethodGet, s.getProjectHandler, true},
		{"/api/v1/project/delete/", http.MethodDelete, s.deleteProjectHandler, true},
		{"/api/v1/auth/logout/{sessionID}", http.MethodDelete, s.logoutUserHandler, true},
		{"/api/v1/deployments", http.MethodGet, s.getAllDeploymentsHandler, true},
		{"/api/v1/deployment", http.MethodGet, s.getDeploymentHandler, true},
		{"/api/v1/project/analytics", http.MethodGet, s.getProjectAnalytics, true},

		{"/api/v1/auth/register", http.MethodPost, s.registerUserHandler, false},
		{"/auth/verify-email", http.MethodGet, s.verifyEmailHandler, false},
		{"/api/v1/create/token", http.MethodPost, s.createVerificationMail, false},
		{"/api/v1/auth/login", http.MethodPost, s.loginUserHandler, false},
		{"/api/v1/auth/refresh", http.MethodPost, s.refreshAccessTokenHandler, false},
		{"/api/v1/user/me", http.MethodGet, s.getUserProfileHandler, true},
	}

	for _, r := range routes {

		handler := Chain(
			http.HandlerFunc(r.handler),
			s.methodMiddleware(r.method),
			s.loggingMiddleware,
		)

		if r.protected {
			handler = Chain(handler, s.authMiddleware)
		}

		mux.Handle(r.path, handler)
	}
}

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {

	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}

func (s *ServerClient) Start() error {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {

		log.Printf("Server running on %s", s.server.Addr)

		if err := s.server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {

			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stop
	log.Println("Gracefully shutting down server...")

	return s.Shutdown()
}

func (s *ServerClient) Shutdown() error {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Server stopped")
	return nil
}
