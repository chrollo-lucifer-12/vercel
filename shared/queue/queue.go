package queue

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

type QueueClient struct {
	client *asynq.Client
}

type EmailJob struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Html    string `json:"html"`
}

type WorkflowJob struct {
	GithubToken string `json:"githubToken"`
	Owner       string `json:"owner"`
	Repo        string `json:"repo"`
	Workflow    string `json:"workflow"`
	Ref         string `json:"ref"`

	GitURL       string `json:"gitURL"`
	ApiURL       string `json:"apiURL"`
	ApiKey       string `json:"apiKey"`
	BucketID     string `json:"bucketId"`
	ProjectSlug  string `json:"projectSlug"`
	DeploymentID string `json:"deploymentId"`
	UserEnv      string `json:"userEnv"`
}

func NewAsynqClient(redisURL string) *QueueClient {
	opt, _ := asynq.ParseRedisURI(redisURL)
	client := asynq.NewClient(opt)
	return &QueueClient{client: client}
}

func (q *QueueClient) NewEmailDeliveryTask(payload EmailJob) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email payload: %w", err)
	}

	task := asynq.NewTask("email:send", data)

	fmt.Println("added new task", task)

	_, err = q.client.Enqueue(task, asynq.Queue("emails"), asynq.MaxRetry(2))

	return task, nil
}

func (q *QueueClient) NewWorkflowTask(payload WorkflowJob) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow payload: %w", err)
	}

	task := asynq.NewTask("workflow:trigger", data)

	fmt.Println("added new workflow task", task)

	_, err = q.client.Enqueue(
		task,
		asynq.Queue("workflows"),
		asynq.MaxRetry(3),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to enqueue workflow task: %w", err)
	}

	return task, nil
}
