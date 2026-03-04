package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chrollo-lucifer-12/api-server/auth"
	"github.com/chrollo-lucifer-12/api-server/server/dto"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/env"
	"github.com/chrollo-lucifer-12/shared/queue"
	"github.com/chrollo-lucifer-12/shared/utils"
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

	dep := &db.Deployment{
		ProjectID: project.ID,
		Status:    "QUEUED",
	}

	err = h.db.CreateDeployment(ctx, dep)
	if err != nil {
		return uuid.Nil, err
	}

	h.queue.NewWorkflowTask(queue.WorkflowJob{
		GithubToken:  env.GithubToken.GetValue(),
		Owner:        "chrollo-lucifer-12",
		Repo:         "vercel",
		Workflow:     "build.yml",
		Ref:          "main",
		GitURL:       project.GitUrl,
		BucketID:     "builds",
		ProjectSlug:  project.SubDomain + strconv.Itoa(dep.Sequence),
		DeploymentID: dep.ID.String(),
		UserEnv:      userEnv,
	})

	return dep.ID, nil
}

func (h *ServerClient) deployHandler(w http.ResponseWriter, r *http.Request) {
	var req DeployRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	project, err := h.db.GetProjectBySlug(ctx, req.ProjectSlug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	cacheKey := "deployments:project:" + project.SubDomain
	cacheKeyProject := fmt.Sprintf("project:slug:%s", project.SubDomain)
	h.redis.Del(ctx, cacheKey)
	h.redis.Del(ctx, cacheKeyProject)

	response := dto.ToCreateDeploymentResponse("queued", project.SubDomain, depID.String())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ServerClient) getAllDeploymentsHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/deployments/")
	if path == "" {
		http.Error(w, "project id required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	cacheKey := "deployments:project:" + path
	cached, err := h.redis.Get(ctx, cacheKey)

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cached))
		return
	}

	projectRes, err := h.db.GetProjectBySlug(ctx, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deployments, err := h.db.GetAllDeployments(ctx, projectRes.ID)

	if err != nil {
		http.Error(w, "failed to get deployments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var d []dto.GetDeploymentResponse
	for _, deployment := range deployments {
		d = append(d, dto.ToGetDeploymentResponse(deployment))
	}

	jsonData, _ := json.Marshal(d)

	h.redis.Set(ctx, cacheKey, jsonData, 30*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(d)
}

func (h *ServerClient) getDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/deployment/")
	if path == "" {
		http.Error(w, "project id required", http.StatusBadRequest)
		return
	}

	deploymentID := utils.StringToUUID(path)

	ctx := r.Context()
	cacheKey := "deployment:" + path

	cached, err := h.redis.Get(ctx, cacheKey)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cached))
		return
	}

	deploymentRes, err := h.db.GetDeploymentByID(ctx, deploymentID)

	if err != nil {
		http.Error(w, "failed to get deployments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.ToGetDeploymentWithLogsResponse(deploymentRes)

	jsonData, _ := json.Marshal(response)
	h.redis.Set(ctx, cacheKey, jsonData, 30*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *ServerClient) getLiveLogs(w http.ResponseWriter, r *http.Request) {

	deploymentId := strings.TrimPrefix(r.URL.Path, "/api/v1/deployment/logs/")
	if deploymentId == "" {
		http.Error(w, "deploymentId required", http.StatusBadRequest)
		return
	}

	streamName := "deployment_logs:" + deploymentId

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	closeNotify := w.(http.CloseNotifier).CloseNotify()
	ctx := r.Context()

	lastID := r.URL.Query().Get("lastID")
	if lastID == "" {
		lastID = "0-0"
	}

	const lastLogs = 50
	messages, _ := h.redis.LRange(ctx, streamName, -lastLogs, -1).Result()
	for _, msg := range messages {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		flusher.Flush()
	}

	for {
		select {
		case <-closeNotify:
			fmt.Println("Client disconnected:", deploymentId)
			return
		default:
			msgs, err := h.redis.StreamRead(ctx, streamName, lastID, 2*time.Second, 10)
			if err != nil {
				fmt.Println("Redis read error:", err)
				time.Sleep(time.Second)
				continue
			}

			if len(msgs) == 0 {
				exists, _ := h.redis.Exists(ctx, streamName)
				if exists == false {
					time.Sleep(5 * time.Second)
					return
				}
				continue
			}

			for _, msg := range msgs {
				lastID = msg.ID
				data := fmt.Sprintf("%v", msg.Values)
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			}
		}
	}
}
