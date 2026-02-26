package server

import (
	"encoding/json"
	"errors"
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

	projects, err := h.db.GetAllProjects(ctx, userID, name, gitURL, limit, offset)
	if err != nil {
		http.Error(w, "error getting projects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var p []dto.CreateProjectResposne
	for _, project := range projects {
		p = append(p, dto.ToCreateProjectResposne(project))
	}

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
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/project/delete/")
	if path == "" {
		http.Error(w, "project id required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	project, err := h.db.GetProjectBySlug(ctx, path)

	query := r.URL.Query()
	fromStr := query.Get("from")
	toStr := query.Get("to")

	var fromTime *time.Time
	var toTime *time.Time

	if fromStr != "" {
		t, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			http.Error(w, "invalid from date", http.StatusBadRequest)
		}
		fromTime = &t
	}

	if toStr != "" {
		t, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			http.Error(w, "invalid from date", http.StatusBadRequest)
		}
		toTime = &t
	}

	analytics, err := h.db.GetAnalytics(ctx, project.SubDomain, fromTime, toTime)

	if err != nil {
		http.Error(w, "failed to get analytics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(analytics)
}
