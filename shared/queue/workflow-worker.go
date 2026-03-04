package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/chrollo-lucifer-12/shared/workflow"
	"github.com/hibiken/asynq"
)

const TypeWorkflowTrigger = "workflow:trigger"

type WorkflowWorker struct {
	server         *asynq.Server
	mux            *asynq.ServeMux
	workflowClient *workflow.WorkflowClient
}

func NewWorkflowWorker(ctx context.Context, token string) *WorkflowWorker {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "127.0.0.1:6379"},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"workflows": 10,
			},
		},
	)

	workflowClient := workflow.NewWorkflowClient(ctx, token)

	mux := asynq.NewServeMux()

	worker := &WorkflowWorker{
		server:         server,
		mux:            mux,
		workflowClient: workflowClient,
	}

	worker.registerHandlers()

	return worker
}

func (w *WorkflowWorker) registerHandlers() {
	w.mux.HandleFunc(TypeWorkflowTrigger, func(ctx context.Context, t *asynq.Task) error {
		var payload WorkflowJob

		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}

		fmt.Println("triggering workflow")

		return w.workflowClient.TriggerWorkflow(ctx, workflow.TriggerWorkflowConfig{
			Owner:        payload.Owner,
			Repo:         payload.Repo,
			WorkflowFile: payload.Workflow,
			Ref:          payload.Ref,
			GithubToken:  payload.GithubToken,
			Inputs: workflow.Input{
				GitURL:       payload.GitURL,
				ApiURL:       payload.ApiURL,
				ApiKey:       payload.ApiKey,
				BucketID:     payload.BucketID,
				ProjectSlug:  payload.ProjectSlug,
				DeploymentID: payload.DeploymentID,
				UserEnv:      payload.UserEnv,
			},
		})
	})
}

func (w *WorkflowWorker) Start() {
	fmt.Println("running workflow worker")
	if err := w.server.Run(w.mux); err != nil {
		log.Fatal(err)
	}
}
