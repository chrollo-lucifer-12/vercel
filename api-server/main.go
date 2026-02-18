package main

import (
	"context"
	"log"

	"github.com/chrollo-lucifer-12/api-server/server"
	"github.com/chrollo-lucifer-12/shared/cache"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/env"
	"github.com/chrollo-lucifer-12/shared/mail"
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

	cacheStore := cache.NewCacheStore(db)

	githubToken := env.GithubToken.GetValue()
	factory := workflow.NewDefaultGithubClientFactory()
	githubClient := factory.NewClient(ctx, githubToken)
	mailClient := mail.NewMailClient(env.ResendApiKey.GetValue())

	validator := workflow.NewDefaultConfigValidator()
	builder := workflow.NewDefaultEventBuilder()

	w := workflow.NewWorkflowClient(
		githubClient,
		cacheStore,
		validator,
		builder,
	)

	h, err := server.NewServerClient(w, db, mailClient)
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := h.Start(); err != nil {
		panic(err)
	}
}
