package main

import (
	"context"
	"log"

	"github.com/chrollo-lucifer-12/api-server/server"
	"github.com/chrollo-lucifer-12/shared/cache"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/env"
	"github.com/chrollo-lucifer-12/shared/workflow"
)

func main() {

	err := env.Load()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	dsn := env.Dsn.GetValue()
	db, err := db.NewDB(dsn, ctx)
	if err != nil {
		panic(err)
	}

	c := cache.NewCacheStore(db)

	githubToken := env.GithubToken.GetValue()
	w := workflow.NewWorkflowClient(ctx, githubToken, c)

	h, err := server.NewServerClient(w, db)
	if err != nil {
		log.Fatal(err)
		return
	}

	h.StartHTTP()
}
