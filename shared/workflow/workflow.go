package workflow

import (
	"context"
	"errors"

	"github.com/google/go-github/v59/github"
	"golang.org/x/oauth2"
)

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

type GithubActionsClient interface {
	CreateWorkflowDispatchEventByFileName(
		ctx context.Context,
		owner string,
		repo string,
		workflowFileName string,
		event github.CreateWorkflowDispatchEventRequest,
	) (*github.Response, error)
}

type MockGithubActionsClient struct {
	CreateWorkflowDispatchEventByFileNameFn func(
		ctx context.Context,
		owner string,
		repo string,
		workflowFileName string,
		event github.CreateWorkflowDispatchEventRequest,
	) (*github.Response, error)

	Called    bool
	LastEvent github.CreateWorkflowDispatchEventRequest
}

type CacheDeleter interface {
	Delete(ctx context.Context, key string) error
}

type ConfigValidator interface {
	Validate(cfg TriggerWorkflowConfig) error
}

type EventBuilder interface {
	Build(cfg TriggerWorkflowConfig) github.CreateWorkflowDispatchEventRequest
}

type WorkflowTrigger interface {
	TriggerWorkflow(ctx context.Context, cfg TriggerWorkflowConfig) error
}

type GithubClientFactory interface {
	NewClient(ctx context.Context, token string) GithubActionsClient
}

type WorkflowClient struct {
	github    GithubActionsClient
	cache     CacheDeleter
	validator ConfigValidator
	builder   EventBuilder
}

func NewWorkflowClient(
	github GithubActionsClient,
	cache CacheDeleter,
	validator ConfigValidator,
	builder EventBuilder,
) *WorkflowClient {
	return &WorkflowClient{
		github:    github,
		cache:     cache,
		validator: validator,
		builder:   builder,
	}
}

func (w *WorkflowClient) TriggerWorkflow(ctx context.Context, cfg TriggerWorkflowConfig) error {
	if err := w.validator.Validate(cfg); err != nil {
		return err
	}

	event := w.builder.Build(cfg)

	if err := w.dispatchWorkflow(ctx, cfg, event); err != nil {
		return err
	}

	return w.cache.Delete(ctx, cfg.Inputs.ProjectSlug)
}

func (w *WorkflowClient) dispatchWorkflow(
	ctx context.Context,
	cfg TriggerWorkflowConfig,
	event github.CreateWorkflowDispatchEventRequest,
) error {
	_, err := w.github.CreateWorkflowDispatchEventByFileName(
		ctx,
		cfg.Owner,
		cfg.Repo,
		cfg.WorkflowFile,
		event,
	)
	return err
}

type DefaultConfigValidator struct{}

func NewDefaultConfigValidator() *DefaultConfigValidator {
	return &DefaultConfigValidator{}
}

func (v *DefaultConfigValidator) Validate(cfg TriggerWorkflowConfig) error {
	if cfg.Owner == "" || cfg.Repo == "" || cfg.WorkflowFile == "" || cfg.Ref == "" {
		return errors.New("invalid workflow config: owner, repo, workflow file, and ref are required")
	}
	if cfg.Inputs.ProjectSlug == "" {
		return errors.New("project slug required")
	}
	return nil
}

type DefaultEventBuilder struct{}

func NewDefaultEventBuilder() *DefaultEventBuilder {
	return &DefaultEventBuilder{}
}

func (b *DefaultEventBuilder) Build(cfg TriggerWorkflowConfig) github.CreateWorkflowDispatchEventRequest {
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

type GithubActionsAdapter struct {
	client *github.Client
}

func NewGithubActionsAdapter(client *github.Client) *GithubActionsAdapter {
	return &GithubActionsAdapter{client: client}
}

func (a *GithubActionsAdapter) CreateWorkflowDispatchEventByFileName(
	ctx context.Context,
	owner string,
	repo string,
	workflowFileName string,
	event github.CreateWorkflowDispatchEventRequest,
) (*github.Response, error) {
	return a.client.Actions.CreateWorkflowDispatchEventByFileName(
		ctx,
		owner,
		repo,
		workflowFileName,
		event,
	)
}

type DefaultGithubClientFactory struct{}

func NewDefaultGithubClientFactory() *DefaultGithubClientFactory {
	return &DefaultGithubClientFactory{}
}

func (f *DefaultGithubClientFactory) NewClient(ctx context.Context, token string) GithubActionsClient {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})
	tc := oauth2.NewClient(ctx, ts)
	return NewGithubActionsAdapter(github.NewClient(tc))
}

type MockCacheDeleter struct {
	DeleteFn func(ctx context.Context, key string) error

	DeletedKey string
	Called     bool
}

func (m *MockCacheDeleter) Delete(ctx context.Context, key string) error {
	m.Called = true
	m.DeletedKey = key

	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, key)
	}

	return nil
}

type MockValidator struct {
	ValidateFn func(cfg TriggerWorkflowConfig) error
}

func (m *MockValidator) Validate(cfg TriggerWorkflowConfig) error {
	if m.ValidateFn != nil {
		return m.ValidateFn(cfg)
	}
	return nil
}

type MockEventBuilder struct {
	BuildFn func(cfg TriggerWorkflowConfig) github.CreateWorkflowDispatchEventRequest

	Called bool
}

func (m *MockEventBuilder) Build(cfg TriggerWorkflowConfig) github.CreateWorkflowDispatchEventRequest {
	m.Called = true

	if m.BuildFn != nil {
		return m.BuildFn(cfg)
	}

	return github.CreateWorkflowDispatchEventRequest{}
}
