package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/env"
	"github.com/chrollo-lucifer-12/shared/workflow"
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

type UpdateHashRequest struct {
	ProjectID string `json:"project_id"`
	Hash      string `json:"hash"`
	GitURL    string `json:"git_url"`
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
			Inputs: workflow.Input{GitURL: project.GitUrl, ApiURL: apiUrl, ApiKey: apiKey, BucketID: "builds", ProjectSlug: project.SubDomain, DeploymentID: dep.ID.String(), UserEnv: userEnv}}); err != nil {
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

	depID, err := h.queueDeployment(ctx, &project, req.UserEnv)
	if err != nil {
		if err == gorm.ErrInvalidData {
			http.Error(w, "another deployment is running", http.StatusConflict)
			return
		}
		http.Error(w, "failed to queue deployment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"status":        "queued",
		"projectSlug":   project.SubDomain,
		"deployment_id": depID.String(),
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

	project := db.Project{
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

	out, err := exec.Command("git", "ls-remote", req.GithubURL, "main").Output()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var hash string
	parts := strings.Split(string(out), "\t")
	if len(parts) > 0 {
		hash = strings.TrimSpace(parts[0])
	}

	err = h.db.CreateHash(ctx, &db.GitHash{ProjectID: project.ID, Hash: hash})
	if err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	//	go h.r.PublishLog(ctx, project.ID.String(), project.GitUrl, hash)

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

	var newLogs []db.LogEvent

	for _, l := range logs {

		deploymentID := l.DeploymentID

		if l.Slug != "" {
			var deployment db.Deployment
			h.db.Raw().Where("slug = ?", l.Slug).Find(&deployment)
			deploymentID = deployment.ID
		}

		if l.Log == "FAILED" || l.Log == "SUCCESS" {
			h.db.UpdateDeployment(ctx, deploymentID, db.Deployment{Status: l.Log})
		}

		newLogs = append(newLogs, db.LogEvent{
			DeploymentID: deploymentID,
			Log:          l.Log,
			Metadata:     l.Metadata,
			Base:         db.Base{CreatedAt: l.CreatedAt},
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

func (h *ServerClient) updateHashHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateHashRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	projectIDParsed, err := uuid.Parse(req.ProjectID)
	if err != nil {
		http.Error(w, "invalid project_id: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	if err := h.db.UpdateHash(ctx, projectIDParsed, db.GitHash{Hash: req.Hash}); err != nil {
		log.Println("update hash error:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	project, err := h.db.GetProjectByID(ctx, projectIDParsed)
	if err != nil {
		http.Error(w, "project not found: "+err.Error(), http.StatusNotFound)
		return
	}

	depID, err := h.queueDeployment(ctx, &project, "")
	if err != nil && err != gorm.ErrInvalidData {
		http.Error(w, "failed to queue deployment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//	go h.r.PublishLog(ctx, req.ProjectID, req.GitURL, req.Hash)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":       "Hash updated and deployment queued",
		"deployment_id": depID.String(),
	})
}
