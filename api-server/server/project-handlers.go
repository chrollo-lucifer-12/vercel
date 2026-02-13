package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/chrollo-lucifer-12/api-server/auth"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/sio/coolname"
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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"name": project.Name,
		"id":   project.ID.String(),
	})
}

func (h *ServerClient) getAllProjectsHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(authKey{}).(*auth.UserClaims)
	userID := claims.ID

	ctx := r.Context()
	projects, err := h.db.GetAllProjects(ctx, userID)

	if err != nil {
		http.Error(w, "error getting projects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(projects)
}

func (h *ServerClient) getProjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	project, _, err := verifyDeployment(r.URL.Path, r, h.db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	projectRes, err := h.db.GetProjectByID(ctx, project.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(projectRes)
}

func (h *ServerClient) deleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	project, _, err := verifyDeployment(r.URL.Path, r, h.db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = h.db.DeleteProject(ctx, project.ID)
	if err != nil {
		http.Error(w, "failed to delete project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ServerClient) getProjectAnalytics(w http.ResponseWriter, r *http.Request) {
	project, _, err := verifyDeployment(r.URL.Path, r, h.db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	ctx := r.Context()
	analytics, err := h.db.GetAnalytics(ctx, project.SubDomain, fromTime, toTime)

	if err != nil {
		http.Error(w, "failed to get analytics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(analytics)
}
