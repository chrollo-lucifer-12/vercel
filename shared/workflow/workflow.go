package workflow

import (
	"context"
	"errors"

	"github.com/chrollo-lucifer-12/shared/cache"
	"github.com/google/go-github/v59/github"
	"golang.org/x/oauth2"
)

type WorkflowClient struct {
	github *github.Client
	cache  CacheDeleter
}

type CacheDeleter interface {
	Delete(ctx context.Context, key string) error
}

func NewWorkflowClient(ctx context.Context, githubToken string, cacheStore *cache.CacheStore) *WorkflowClient {
	return &WorkflowClient{
		github: newGithubClient(ctx, githubToken),
		cache:  cacheStore,
	}
}

func newGithubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})

	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
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

func (w *WorkflowClient) TriggerWorkflow(ctx context.Context, cfg TriggerWorkflowConfig) error {

	if err := validateConfig(cfg); err != nil {
		return err
	}

	event := buildWorkflowDispatchEvent(cfg)

	if err := w.dispatchWorkflow(ctx, cfg, event); err != nil {
		return err
	}

	return w.invalidateCache(ctx, cfg.Inputs.ProjectSlug)
}

func validateConfig(cfg TriggerWorkflowConfig) error {
	if cfg.Owner == "" ||
		cfg.Repo == "" ||
		cfg.WorkflowFile == "" ||
		cfg.Ref == "" {
		return errors.New("invalid workflow config")
	}

	if cfg.Inputs.ProjectSlug == "" {
		return errors.New("project slug required")
	}

	return nil
}

func buildWorkflowDispatchEvent(cfg TriggerWorkflowConfig) github.CreateWorkflowDispatchEventRequest {

	return github.CreateWorkflowDispatchEventRequest{
		Ref: cfg.Ref,
		Inputs: map[string]interface{}{
			"gitURL":       cfg.Inputs.GitURL,
			"apiURL":       cfg.Inputs.ApiURL,
			"apiKey":       cfg.Inputs.ApiKey,
			"bucketId":     cfg.Inputs.BucketID,
			"projectSlug":  cfg.Inputs.ProjectSlug,
			"deploymentId": cfg.Inputs.DeploymentID,
			"userEnv":      cfg.Inputs.UserEnv,
		},
	}
}

func (w *WorkflowClient) dispatchWorkflow(
	ctx context.Context,
	cfg TriggerWorkflowConfig,
	event github.CreateWorkflowDispatchEventRequest,
) error {

	_, err := w.github.Actions.CreateWorkflowDispatchEventByFileName(
		ctx,
		cfg.Owner,
		cfg.Repo,
		cfg.WorkflowFile,
		event,
	)

	return err
}

func (w *WorkflowClient) invalidateCache(ctx context.Context, projectSlug string) error {
	return w.cache.Delete(ctx, projectSlug)
}
