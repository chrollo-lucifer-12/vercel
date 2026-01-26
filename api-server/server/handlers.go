package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/chrollo-lucifer-12/api-server/models"
	"github.com/google/uuid"
	"github.com/sio/coolname"
	"gorm.io/gorm"
)

type DeployRequest struct {
	ProjectID string `json:"project_id"`
}

type ProjectRequest struct {
	ProjectName string `json:"project_name"`
	GithubURL   string `json:"github_url"`
}

func (h *ServerClient) deployHandler(w http.ResponseWriter, r *http.Request) {
	var req DeployRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.ProjectID == "" {
		http.Error(w, "project_id is required", http.StatusBadRequest)
		return
	}

	projectIDUUID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		http.Error(w, "invalid project_id: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	project, err := h.db.GetProjectByID(ctx, projectIDUUID)
	if err != nil {
		http.Error(w, "project not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if project.SubDomain == "" || project.GitUrl == "" {
		http.Error(w, "project data is incomplete", http.StatusInternalServerError)
		return
	}

	d, err := gorm.G[models.Deployment](h.db.Raw()).Where("project_id = ?", project.ID).Where("status IN ?", []string{"QUEUED", "PENDING"}).Find(ctx)
	if err != nil {
		http.Error(w, "failed to fetch deployment "+err.Error(), http.StatusNotFound)
		return
	}

	if len(d) > 0 {
		http.Error(w, "another deployment is running", http.StatusConflict)
		return
	}

	tx := h.db.Raw().Begin()
	if tx.Error != nil {
		http.Error(w, "failed to start transaction", http.StatusInternalServerError)
		return
	}

	dep := &models.Deployment{
		ProjectID: projectIDUUID,
		Status:    "QUEUED",
	}

	if err := tx.Create(dep).Error; err != nil {
		tx.Rollback()
		http.Error(w, "failed to create deployment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		http.Error(w, "failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	projectSlug := project.SubDomain
	gitUrl := project.GitUrl

	if err := h.wClient.TriggerWorkflow(ctx, gitUrl, projectSlug, dep.ID.String()); err != nil {

		_ = h.db.Raw().Model(&models.Deployment{}).
			Where("id = ?", dep.ID).
			Update("status", "FAILED")

		http.Error(w, "failed to trigger workflow: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"status":      "queued",
		"projectSlug": projectSlug,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ServerClient) projectHandler(w http.ResponseWriter, r *http.Request) {
	var req ProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	subdomain, err := coolname.Slug()
	if err != nil {
		http.Error(w, "failed to generate subdomain", http.StatusInternalServerError)
		return
	}

	project := models.Project{
		Name:      req.ProjectName,
		GitUrl:    req.GithubURL,
		SubDomain: subdomain,
	}

	ctx := context.Background()
	if err := h.db.CreateProject(ctx, &project); err != nil {
		log.Println("create project error:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"name": project.Name,
		"id":   project.ID.String(),
	})
}

func (h *ServerClient) logsHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 3 || parts[1] != "logs" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	deploymentID := parts[2]
	if deploymentID == "" {
		http.Error(w, "deployment_id is required", http.StatusBadRequest)
		return
	}

	logs, err := h.clickDB.GetLogsByDeployment(r.Context(), deploymentID)
	if err != nil {
		http.Error(w, "failed to fetch logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
