package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/chrollo-lucider-12/background/redis"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisURL := os.Getenv("REDIS_URL")
	apiURL := os.Getenv("API_URL")

	r, err := redis.NewRedisClient(redisURL, apiURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1)

	go r.SubscribeProxyLogs(ctx, "analytics_stream")
	go r.SubscribeStreams(ctx, "logs_stream")
	go r.SubscribeHashStreams(ctx, "hash_stream")

	wg.Wait()
}
