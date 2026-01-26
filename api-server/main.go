package main

import (
	"context"
	"log"

	"github.com/chrollo-lucifer-12/api-server/clickhouse"
	"github.com/chrollo-lucifer-12/api-server/env"
	"github.com/chrollo-lucifer-12/api-server/models"
	"github.com/chrollo-lucifer-12/api-server/redis"
	"github.com/chrollo-lucifer-12/api-server/server"
	"github.com/chrollo-lucifer-12/api-server/workflow"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	e, err := env.NewEnv()
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx := context.Background()

	c := clickhouse.NewClickHouseDB(e.CLICKHOUSE_ADDR, e.CLICKHOUSE_USERNAME, e.CLICKHOUSE_PASSWORD)

	r, err := redis.NewRedisClient(e.REDIS_URL, c)
	if err != nil {
		log.Fatal(err)
		return
	}

	go r.SubscribeStreams(ctx, "logs_stream")

	db, err := models.NewDB(e.DSN, ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	w := workflow.NewWorkflowClient(ctx)

	h, err := server.NewServerClient(w, db)
	if err != nil {
		log.Fatal(err)
		return
	}

	h.StartHTTP()
}
