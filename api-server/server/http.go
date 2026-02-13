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
	"github.com/chrollo-lucifer-12/shared/workflow"
)

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

	server := &ServerClient{
		wClient: wClient,
		db:      db,
		auth:    auth,
	}

	mux := http.NewServeMux()
	server.registerRoutes(mux)

	server.server = &http.Server{
		Addr:    ":9000",
		Handler: mux,
	}

	return &ServerClient{wClient: wClient, db: db, auth: auth}, nil
}

func (s *ServerClient) registerRoutes(mux *http.ServeMux) {

	protectedRoutes := map[string]http.HandlerFunc{
		"/api/v1/deploy/create":    s.deployHandler,
		"/api/v1/project/create":   s.createProjectHandler,
		"/api/v1/projects":         s.getAllProjectsHandler,
		"/api/v1/project":          s.getProjectHandler,
		"/api/v1/project/delete":   s.deleteProjectHandler,
		"/api/v1/auth/logout":      s.logoutUserHandler,
		"/api/v1/deployments":      s.getAllDeploymentsHandler,
		"/api/v1/deplyoment":       s.getDeploymentHandler,
		"api/v1/project/analytics": s.getProjectAnalytics,
	}

	for path, handler := range protectedRoutes {
		mux.Handle(path, Chain(handler, s.authMiddleware))
	}

	mux.HandleFunc("/api/v1/auth/register", s.registerUserHandler)
	mux.HandleFunc("/api/v1/auth/login", s.loginUserHandler)
	mux.HandleFunc("/api/v1/auth/refresh", s.refreshAccessTokenHandler)

}

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func (h *ServerClient) Start() error {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s", h.server.Addr)
		if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server gracefully...")

	return h.Shutdown()
}

func (h *ServerClient) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Server stopped")
	return nil
}
