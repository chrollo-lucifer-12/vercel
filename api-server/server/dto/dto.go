package dto

import (
	"time"

	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
}

type LoginResponse struct {
	SessionID             string    `json:"session_id"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	AccessToken           string    `json:"access_token"`
	User                  UserResponse
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	SessionID            string    `json:"session_id"`
}

type CreateDeploymentResponse struct {
	Status       string `json:"status"`
	ProjectSlug  string `json:"project_slug"`
	DeploymentID string `json:"deployment_id"`
}

type GetDeploymentResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
	Sequence  int       `json:"sequence"`
}

type LogsResponse struct {
	Log       string         `json:"log"`
	CreatedAt time.Time      `json:"created_at"`
	Metadata  datatypes.JSON `json:"metadata"`
}

type GetDeploymentWithLogsResponse struct {
	Deployment GetDeploymentResponse
	Logs       []LogsResponse
}

type CreateProjectResposne struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	SubDomain string    `json:"sub_domain"`
	CreatedAt time.Time `json:"created_at"`
	GitUrl    string    `json:"git_url"`
}

type GetProjectWithDeployment struct {
	Project    CreateProjectResposne
	Deployment GetDeploymentWithLogsResponse
}

func ToGetProjectWithDeployment(project db.Project, deployment db.Deployment) GetProjectWithDeployment {
	return GetProjectWithDeployment{
		Project:    ToCreateProjectResposne(project),
		Deployment: ToGetDeploymentWithLogsResponse(deployment),
	}
}

func ToCreateProjectResposne(project db.Project) CreateProjectResposne {
	return CreateProjectResposne{
		ID:        project.ID.String(),
		Name:      project.Name,
		SubDomain: project.SubDomain,
		CreatedAt: project.CreatedAt,
		GitUrl:    project.GitUrl,
	}
}

func ToLogsResponse(logs []db.LogEvent) []LogsResponse {
	var r []LogsResponse

	for _, log := range logs {
		r = append(r, LogsResponse{Log: log.Log, CreatedAt: log.CreatedAt, Metadata: log.Metadata})
	}

	return r
}

func ToGetDeploymentWithLogsResponse(deployment db.Deployment) GetDeploymentWithLogsResponse {
	return GetDeploymentWithLogsResponse{
		Deployment: ToGetDeploymentResponse(deployment),
		Logs:       ToLogsResponse(deployment.LogEvents),
	}
}

func ToGetDeploymentResponse(deployment db.Deployment) GetDeploymentResponse {
	return GetDeploymentResponse{
		ID:        deployment.ID,
		CreatedAt: deployment.CreatedAt,
		Status:    deployment.Status,
		Sequence:  deployment.Sequence,
	}
}

func ToCreateDeploymentResponse(status, projectSlug, deploymentID string) CreateDeploymentResponse {
	return CreateDeploymentResponse{
		Status:       status,
		ProjectSlug:  projectSlug,
		DeploymentID: deploymentID,
	}
}

func ToRenewAccessTokenResponse(accessToken string, accessTokenExpiresAt time.Time, sessionID string) RenewAccessTokenResponse {
	return RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenExpiresAt,
		SessionID:            sessionID,
	}
}

func ToUserResponse(user db.User) UserResponse {
	return UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		IsVerified: user.IsVerified,
	}
}

func ToLoginResponse(user db.User, sessionID string, refreshToken string, accessTokenExpiresAt time.Time, refreshTokenExpiresAt time.Time, accessToken string) LoginResponse {
	return LoginResponse{
		SessionID:             sessionID,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
		AccessToken:           accessToken,
		User:                  ToUserResponse(user),
	}
}
