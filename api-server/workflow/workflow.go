package workflow

import (
	"context"
	"os"

	"github.com/chrollo-lucifer-12/api-server/redis"
	"github.com/google/go-github/v59/github"
	"golang.org/x/oauth2"
)

type WorkflowClient struct {
	wClient *github.Client
	rClient *redis.RedisClient
}

func NewWorkflowClient(ctx context.Context, rClient *redis.RedisClient) *WorkflowClient {
	token := os.Getenv("GITHUB_TOKEN")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &WorkflowClient{wClient: client, rClient: rClient}
}

func (w *WorkflowClient) TriggerWorkflow(ctx context.Context, gitURL, projectSlug string) error {
	apiURL := os.Getenv("API_URL")
	apiKey := os.Getenv("API_KEY")
	redisURL := os.Getenv("REDIS_URL")
	owner := "chrollo-lucifer-12"
	repo := "vercel"
	workflowFile := "build.yml"

	go w.rClient.SubscribeChannel(ctx, "logs:"+projectSlug)

	inputs := map[string]interface{}{
		"gitURL":      gitURL,
		"apiURL":      apiURL,
		"apiKey":      apiKey,
		"bucketId":    "builds",
		"projectSlug": projectSlug,
		"redisURL":    redisURL,
	}

	event := github.CreateWorkflowDispatchEventRequest{
		Ref:    "main",
		Inputs: inputs,
	}

	_, err := w.wClient.Actions.CreateWorkflowDispatchEventByFileName(ctx, owner, repo, workflowFile, event)

	return err
}
