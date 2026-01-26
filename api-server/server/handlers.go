package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/chrollo-lucifer-12/api-server/models"
	"github.com/google/uuid"
	"github.com/sio/coolname"
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

	projectSlug := project.SubDomain
	gitUrl := project.GitUrl

	if err := h.wClient.TriggerWorkflow(ctx, gitUrl, projectSlug); err != nil {
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
