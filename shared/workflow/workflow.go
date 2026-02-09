package workflow

import (
	"context"

	"github.com/chrollo-lucifer-12/shared/cache"
	"github.com/google/go-github/v59/github"
	"golang.org/x/oauth2"
)

type WorkflowClient struct {
	wClient *github.Client
	cache   *cache.CacheStore
}

func NewWorkflowClient(ctx context.Context, githubToken string, cache *cache.CacheStore) *WorkflowClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &WorkflowClient{wClient: client, cache: cache}
}

type Input struct {
	GitURL       string
	ApiURL       string
	ApiKey       string
	BucketID     string
	ProjectSlug  string
	DeploymentID string
	UserEnv      string
}

type TriggerWorkflowConfig struct {
	Inputs       Input
	Owner        string
	Repo         string
	WorkflowFile string
	Ref          string
}

func (w *WorkflowClient) TriggerWorkflow(ctx context.Context, config TriggerWorkflowConfig) error {

	owner := config.Owner
	repo := config.Repo
	workflowFile := config.WorkflowFile

	inputs := map[string]interface{}{
		"gitURL":       config.Inputs.GitURL,
		"apiURL":       config.Inputs.ApiURL,
		"apiKey":       config.Inputs.ApiKey,
		"bucketId":     config.Inputs.BucketID,
		"projectSlug":  config.Inputs.ProjectSlug,
		"deploymentId": config.Inputs.DeploymentID,
		"userEnv":      config.Inputs.UserEnv,
	}

	event := github.CreateWorkflowDispatchEventRequest{
		Ref:    config.Ref,
		Inputs: inputs,
	}

	_, err := w.wClient.Actions.CreateWorkflowDispatchEventByFileName(ctx, owner, repo, workflowFile, event)
	if err != nil {
		return err
	}
	err = w.cache.Delete(ctx, config.Inputs.ProjectSlug)

	return err
}
