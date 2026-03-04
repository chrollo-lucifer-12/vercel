package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chrollo-lucifer-12/api-server/auth"
	"github.com/chrollo-lucifer-12/api-server/server/dto"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/google/uuid"
	"github.com/sio/coolname"
	"gorm.io/gorm"
)

func (h *ServerClient) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	var req ProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(authKey{}).(*auth.UserClaims)
	userID := claims.ID

	subdomain, err := coolname.Slug()
	if err != nil {
		http.Error(w, "failed to generate subdomain", http.StatusInternalServerError)
		return
	}

	project := db.Project{
		Name:      req.ProjectName,
		GitUrl:    req.GithubURL,
		SubDomain: subdomain,
		UserID:    userID,
	}

	ctx := r.Context()
	if err := h.db.CreateProject(ctx, &project); err != nil {
		log.Println("create project error:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	pattern := fmt.Sprintf("projects:user:%s:*", userID)
	if err := h.redis.DeleteByPattern(ctx, pattern); err != nil {
		log.Println("cache invalidation error:", err)
	}

	response := dto.ToCreateProjectResposne(project)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *ServerClient) getAllProjectsHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(authKey{}).(*auth.UserClaims)
	userID := claims.ID
	ctx := r.Context()

	query := r.URL.Query()

	name := query.Get("name")
	gitURL := query.Get("giturl")

	limit := 10
	offset := 0

	if l := query.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := query.Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	cacheKey := fmt.Sprintf(
		"projects:user:%s:name:%s:git:%s:limit:%d:offset:%d",
		userID,
		name,
		gitURL,
		limit,
		offset,
	)

	cached, err := h.redis.Get(ctx, cacheKey)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cached))
		return
	}

	projects, err := h.db.GetAllProjects(ctx, userID, name, gitURL, limit, offset)
	if err != nil {
		http.Error(w, "error getting projects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var p []dto.CreateProjectResposne
	for _, project := range projects {
		p = append(p, dto.ToCreateProjectResposne(project))
	}

	jsonData, _ := json.Marshal(p)
	h.redis.Set(ctx, cacheKey, jsonData, 30*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(p)
}

func (h *ServerClient) getProjectHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/project/")
	if path == "" {
		http.Error(w, "project id required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	cacheKey := fmt.Sprintf("project:slug:%s", path)

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

	deployment, err := h.db.GetLatestDeployment(ctx, projectRes.ID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.ToGetProjectWithDeployment(projectRes, deployment)

	jsonData, _ := json.Marshal(response)

	h.redis.Set(ctx, cacheKey, jsonData, 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func (h *ServerClient) deleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/project/delete/")
	if path == "" {
		http.Error(w, "project id required", http.StatusBadRequest)
		return
	}

	projectID, err := uuid.Parse(path)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	err = h.db.DeleteProject(ctx, projectID)
	if err != nil {
		http.Error(w, "failed to delete project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ServerClient) getProjectAnalytics(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/project/analytics/")
	if path == "" {
		http.Error(w, "project id required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	query := r.URL.Query()
	fromStr := query.Get("from")
	toStr := query.Get("to")

	var fromTime *time.Time
	var toTime *time.Time

	if fromStr != "" {
		t, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			http.Error(w, "invalid from date", http.StatusBadRequest)
			return
		}
		fromTime = &t
	}

	if toStr != "" {
		t, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			http.Error(w, "invalid to date", http.StatusBadRequest)
			return
		}
		toTime = &t
	}

	fromKey := "nil"
	toKey := "nil"

	if fromTime != nil {
		fromKey = fromTime.Format("2006-01-02")
	}
	if toTime != nil {
		toKey = toTime.Format("2006-01-02")
	}

	cacheKey := fmt.Sprintf(
		"analytics:slug:%s:from:%s:to:%s",
		path,
		fromKey,
		toKey,
	)

	cached, err := h.redis.Get(ctx, cacheKey)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cached))
		return
	}

	project, err := h.db.GetProjectBySlug(ctx, path)
	if err != nil {
		http.Error(w, "project not found", http.StatusInternalServerError)
		return
	}

	analytics, err := h.db.GetAnalytics(ctx, project.SubDomain, fromTime, toTime)
	if err != nil {
		http.Error(w, "failed to get analytics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, _ := json.Marshal(analytics)

	h.redis.Set(ctx, cacheKey, jsonData, 2*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
