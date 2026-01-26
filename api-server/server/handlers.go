package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sio/coolname"
)

type ProjectRequest struct {
	GithubURL string `json:"github_url"`
}

func (h *ServerClient) projectHandler(w http.ResponseWriter, r *http.Request) {
	var req ProjectRequest
	json.NewDecoder(r.Body).Decode(&req)

	projectSlug, _ := coolname.Slug()

	gitUrl := req.GithubURL

	ctx := context.Background()

	err := h.wClient.TriggerWorkflow(ctx, gitUrl, projectSlug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"status":      "queued",
		"projectSlug": projectSlug,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
