package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/chrollo-lucifer-12/api-server/models"
	"github.com/google/uuid"
	"github.com/sio/coolname"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DeployRequest struct {
	ProjectID string `json:"project_id"`
	UserEnv   string `json:"user_env"`
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

	if err := h.wClient.TriggerWorkflow(ctx, gitUrl, projectSlug, dep.ID.String(), req.UserEnv); err != nil {

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

func (h *ServerClient) registerLogsRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var logs []LogRequest
	if err := json.NewDecoder(r.Body).Decode(&logs); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	var newLogs []models.LogEvent

	for _, l := range logs {

		deploymentID := l.DeploymentID

		if l.Slug != "" {
			var deployment models.Deployment
			h.db.Raw().Where("slug = ?", l.Slug).Find(&deployment)
			deploymentID = deployment.ID
		}

		newLogs = append(newLogs, models.LogEvent{
			DeploymentID: deploymentID,
			Log:          l.Log,
			Metadata:     l.Metadata,
			Base:         models.Base{CreatedAt: l.CreatedAt},
		})

	}

	err := h.db.CreateLogEvents(ctx, &newLogs)
	if err != nil {
		log.Println("failed to insert log:", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *ServerClient) logsHandler(w http.ResponseWriter, r *http.Request) {

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	deploymentID := r.URL.Query().Get("deployment_id")
	if deploymentID == "" {
		http.Error(w, "deployment_id is required", http.StatusBadRequest)
		return
	}
	deploymentIDParsed, _ := uuid.Parse(deploymentID)
	var fromTime, toTime time.Time
	var err error

	if from != "" {
		fromTime, err = time.Parse(time.RFC3339, from)
		if err != nil {
			http.Error(w, "invalid 'from' date format (use RFC3339)", http.StatusBadRequest)
			return
		}
	}

	if to != "" {
		toTime, err = time.Parse(time.RFC3339, to)
		if err != nil {
			http.Error(w, "invalid 'to' date format (use RFC3339)", http.StatusBadRequest)
			return
		}
	}

	logs, err := h.db.GetLogEventsByDeploymentAndTimeRange(context.Background(), deploymentIDParsed, fromTime, toTime)
	if err != nil {
		http.Error(w, "failed to fetch logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *ServerClient) analyticsHandler(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	pagePath := queryParams.Get("path")
	statusCode := queryParams.Get("status_code")
	deploymentID := queryParams.Get("deployment_id")
	if deploymentID == "" {
		http.Error(w, "deployment_id is required", http.StatusBadRequest)
		return
	}
	deploymentIDParsed, _ := uuid.Parse(deploymentID)

	var fromTime, toTime time.Time
	if from := queryParams.Get("from"); from != "" {
		t, err := time.Parse("2006-01-02", from)
		if err != nil {
			http.Error(w, "invalid from date (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		fromTime = t
	}

	if to := queryParams.Get("to"); to != "" {
		t, err := time.Parse("2006-01-02", to)
		if err != nil {
			http.Error(w, "invalid to date (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		toTime = t
	}

	logs, err := h.db.GetAnalytics(context.Background(), deploymentIDParsed, fromTime, toTime, statusCode, pagePath)
	if err != nil {
		http.Error(w, "failed to fetch logs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
