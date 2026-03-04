package workflow

import (
	"context"
	"fmt"

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
	GithubToken  string
}

type WorkflowClient struct {
	client *github.Client
}

func NewWorkflowClient(ctx context.Context, token string) *WorkflowClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &WorkflowClient{
		client: github.NewClient(tc),
	}
}

func (w *WorkflowClient) TriggerWorkflow(
	ctx context.Context,
	cfg TriggerWorkflowConfig,
) error {

	event := github.CreateWorkflowDispatchEventRequest{
		Ref: cfg.Ref,
		Inputs: map[string]interface{}{
			"gitURL":       cfg.Inputs.GitURL,
			"bucketId":     cfg.Inputs.BucketID,
			"projectSlug":  cfg.Inputs.ProjectSlug,
			"deploymentId": cfg.Inputs.DeploymentID,
			"userEnv":      "test",
		},
	}

	_, err := w.client.Actions.CreateWorkflowDispatchEventByFileName(
		ctx,
		cfg.Owner,
		cfg.Repo,
		cfg.WorkflowFile,
		event,
	)

	fmt.Println(err)

	return err
}
