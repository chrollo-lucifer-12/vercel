package main

import (
	"context"
	"sync"

	"github.com/chrollo-lucifer-12/shared/env"
	"github.com/chrollo-lucifer-12/shared/queue"
)

func main() {
	env.Load()
	ctx := context.TODO()

	emailWorker := queue.NewEmailWorkerServer(env.RedisUrl.GetValue(), env.ResendApiKey.GetValue())
	workflowWorker := queue.NewWorkflowWorker(ctx, env.GithubToken.GetValue(), env.RedisUrl.GetValue())
	analyticsWorker := queue.NewAnalyticsWorker(ctx, env.Dsn.GetValue(), env.RedisUrl.GetValue())

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		emailWorker.Start()
	}()

	go func() {
		defer wg.Done()
		workflowWorker.Start()
	}()

	go func() {
		defer wg.Done()
		analyticsWorker.Start()
	}()

	wg.Wait()
}
