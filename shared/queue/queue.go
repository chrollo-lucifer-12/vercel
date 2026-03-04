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
