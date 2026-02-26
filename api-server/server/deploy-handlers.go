package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/chrollo-lucifer-12/api-server/auth"
	"github.com/chrollo-lucifer-12/api-server/server/dto"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/env"
	"github.com/chrollo-lucifer-12/shared/workflow"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func verifyDeployment(path string, r *http.Request, db *db.DB) (*db.Project, *db.Deployment, error) {
	parts := strings.Split(path, "/")

	if len(parts) < 2 || parts[1] == "" {
		return nil, nil, fmt.Errorf("invalid path")
	}

	projectID, err := uuid.Parse(parts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid project id")
	}

	claims, ok := r.Context().Value(authKey{}).(*auth.UserClaims)
	if !ok {
		return nil, nil, fmt.Errorf("unauthorized")
	}

	ctx := r.Context()

	project, err := db.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, nil, fmt.Errorf("project not found: %w", err)
	}

	if claims.ID != project.UserID {
		return nil, nil, fmt.Errorf("forbidden")
	}

	if len(parts) < 3 || parts[2] == "" {
		return &project, nil, nil
	}

	deploymentID, err := uuid.Parse(parts[2])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid deployment id")
	}

	deployment, err := db.GetDeploymentByID(ctx, deploymentID)
	if err != nil {
		return nil, nil, fmt.Errorf("deployment not found: %w", err)
	}

	if deployment.ProjectID != project.ID {
		return nil, nil, fmt.Errorf("deployment does not belong to project")
	}

	return &project, &deployment, nil
}

func (h *ServerClient) queueDeployment(ctx context.Context, project *db.Project, userEnv string) (uuid.UUID, error) {
	apiUrl := env.SupabaseUrl.GetValue()
	apiKey := env.SupabaseSecret.GetValue()
	d, err := gorm.G[db.Deployment](h.db.Raw()).
		Where("project_id = ?", project.ID).
		Where("status IN ?", []string{"QUEUED", "PENDING"}).
		Find(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	if len(d) > 0 {
		return uuid.Nil, gorm.ErrInvalidData
	}

	tx := h.db.Raw().Begin()
	if tx.Error != nil {
		return uuid.Nil, tx.Error
	}

	dep := &db.Deployment{
		ProjectID: project.ID,
		Status:    "QUEUED",
	}

	if err := tx.Create(dep).Error; err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	go func() {
		if err := h.wClient.TriggerWorkflow(ctx, workflow.TriggerWorkflowConfig{Owner: "chrollo-lucifer-12", Repo: "vercel", WorkflowFile: "build.yml", Ref: "main",
			Inputs: workflow.Input{GitURL: project.GitUrl, ApiURL: apiUrl, ApiKey: apiKey, BucketID: "builds", ProjectSlug: project.SubDomain + strconv.Itoa(dep.Sequence), DeploymentID: dep.ID.String(), UserEnv: userEnv}}); err != nil {
			_ = h.db.Raw().Model(&db.Deployment{}).
				Where("id = ?", dep.ID).
				Update("status", "FAILED")
			log.Println("failed to trigger workflow:", err)
		} else {
			_ = h.db.Raw().Model(&db.Deployment{}).
				Where("id = ?", dep.ID).
				Update("status", "PENDING")
		}
	}()

	return dep.ID, nil
}

func (h *ServerClient) deployHandler(w http.ResponseWriter, r *http.Request) {
	var req DeployRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	project, _, err := verifyDeployment(r.URL.Path, r, h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	depID, err := h.queueDeployment(ctx, project, req.UserEnv)
	if err != nil {
		if err == gorm.ErrInvalidData {
			http.Error(w, "another deployment is running", http.StatusConflict)
			return
		}
		http.Error(w, "failed to queue deployment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.ToCreateDeploymentResponse("queued", project.SubDomain, depID.String())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ServerClient) getAllDeploymentsHandler(w http.ResponseWriter, r *http.Request) {
	project, _, err := verifyDeployment(r.URL.Path, r, h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	deployments, err := h.db.GetAllDeployments(ctx, project.ID)

	if err != nil {
		http.Error(w, "failed to get deployments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var d []dto.GetDeploymentResponse
	for _, deployment := range deployments {
		d = append(d, dto.ToGetDeploymentResponse(deployment))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(d)
}

func (h *ServerClient) getDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	_, deployment, err := verifyDeployment(r.URL.Path, r, h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	deploymentRes, err := h.db.GetDeploymentByID(ctx, deployment.ID)

	if err != nil {
		http.Error(w, "failed to get deployments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.ToGetDeploymentWithLogsResponse(deploymentRes)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
