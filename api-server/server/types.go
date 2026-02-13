package server

import (
	"net/http"
	"time"

	"github.com/chrollo-lucifer-12/api-server/auth"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/workflow"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type DeployRequest struct {
	UserEnv string `json:"user_env"`
}

type ProjectRequest struct {
	ProjectName string `json:"project_name"`
	GithubURL   string `json:"github_url"`
}

type LogRequest struct {
	DeploymentID uuid.UUID      `json:"deployment_id"`
	Log          string         `json:"log"`
	Metadata     datatypes.JSON `json:"metadata"`
	CreatedAt    time.Time      `json:"created_at"`
	Slug         string         `json:"slug"`
}

type UpdateHashRequest struct {
	ProjectID string `json:"project_id"`
	Hash      string `json:"hash"`
	GitURL    string `json:"git_url"`
}

type authKey struct{}

type Middleware func(http.Handler) http.Handler

type ServerClient struct {
	wClient *workflow.WorkflowClient
	db      *db.DB
	auth    *auth.AuthService
	server  *http.Server
}

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
	SessionID             string    `json:"session_id"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	AccessToken           string    `json:"access_token"`
	User                  UserRes
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RenewAccessTokenRes struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
