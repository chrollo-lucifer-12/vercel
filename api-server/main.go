package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/v59/github"
	"github.com/sio/coolname"
	"golang.org/x/oauth2"
)

type ProjectRequest struct {
	GithubURL string `json:"github_url"`
}

func projectHandler(w http.ResponseWriter, r *http.Request) {
	var req ProjectRequest
	json.NewDecoder(r.Body).Decode(&req)

	projectSlug, _ := coolname.Slug()
	gitUrl := req.GithubURL

	err := triggerWorkflow(gitUrl, projectSlug)
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

func triggerWorkflow(gitURL, projectSlug string) error {
	ctx := context.Background()
	token := os.Getenv("GITHUB_TOKEN")
	apiURL := os.Getenv("API_URL")
	apiKey := os.Getenv("API_KEY")

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	owner := "chrollo-lucifer-12"
	repo := "vercel"
	workflowFile := "build.yml"

	inputs := map[string]interface{}{
		"gitURL":   gitURL,
		"apiURL":   apiURL,
		"apiKey":   apiKey,
		"buckerId": "builds",
	}

	event := github.CreateWorkflowDispatchEventRequest{
		Ref:    "main",
		Inputs: inputs,
	}

	_, err := client.Actions.CreateWorkflowDispatchEventByFileName(ctx, owner, repo, workflowFile, event)
	return err
}

func main() {
	http.HandleFunc("/project", projectHandler)

	log.Fatal(http.ListenAndServe(":9000", nil))
}
