package main

import (
	"context"
	"log"

	"github.com/chrollo-lucifer-12/api-server/server"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/env"
	"github.com/chrollo-lucifer-12/shared/queue"
	"github.com/chrollo-lucifer-12/shared/redis"
	"github.com/chrollo-lucifer-12/shared/workflow"
)

func main() {

	env.Load()

	ctx := context.Background()

	dsn := env.Dsn.GetValue()
	db, err := db.NewDB(dsn, ctx)
	if err != nil {
		panic(err)
	}
	githubToken := env.GithubToken.GetValue()
	factory := workflow.NewDefaultGithubClientFactory()
	githubClient := factory.NewClient(ctx, githubToken)
	redisClient := redis.NewRedisClient(env.RedisUrl.GetValue())
	queueClient := queue.NewAsynqClient(env.RedisUrl.GetValue())

	validator := workflow.NewDefaultConfigValidator()
	builder := workflow.NewDefaultEventBuilder()

	w := workflow.NewWorkflowClient(
		githubClient,
		validator,
		builder,
	)

	h, err := server.NewServerClient(w, db, redisClient, queueClient)
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := h.Start(); err != nil {
		panic(err)
	}
}
