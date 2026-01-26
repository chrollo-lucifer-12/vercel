package main

import (
	"context"
	"log"
	"os"

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
	}

	redisURL := os.Getenv("REDIS_URL")
	dsn := os.Getenv("DSN")
	ctx := context.Background()

	_, err = models.NewDB(dsn, ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	r, err := redis.NewRedisClient(redisURL)

	if err != nil {
		log.Fatal(err)
		return
	}

	w := workflow.NewWorkflowClient(ctx, r)

	h, err := server.NewServerClient(w)
	if err != nil {
		log.Fatal(err)
		return
	}

	h.StartHTTP()
}
