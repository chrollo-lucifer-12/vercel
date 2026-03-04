package main

import (
	"context"
	"log"

	"github.com/chrollo-lucifer-12/api-server/server"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/env"
	"github.com/chrollo-lucifer-12/shared/queue"
	"github.com/chrollo-lucifer-12/shared/redis"
)

func main() {

	env.Load()

	ctx := context.Background()

	dsn := env.Dsn.GetValue()
	db, err := db.NewDB(dsn, ctx)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewRedisClient(env.RedisUrl.GetValue())
	queueClient := queue.NewAsynqClient(env.RedisUrl.GetValue())

	h, err := server.NewServerClient(db, redisClient, queueClient)
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := h.Start(); err != nil {
		panic(err)
	}
}
