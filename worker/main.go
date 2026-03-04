package main

import (
	"context"

	"github.com/chrollo-lucifer-12/shared/env"
	"github.com/chrollo-lucifer-12/shared/queue"
)

func main() {
	env.Load()
	ctx := context.TODO()

	emailWorker := queue.NewEmailWorkerServer(env.RedisUrl.GetValue(), env.ResendApiKey.GetValue())
	workflowWorker := queue.NewWorkflowWorker(ctx, env.GithubToken.GetValue())

	go emailWorker.Start()
	go workflowWorker.Start()
}
