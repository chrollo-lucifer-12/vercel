package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	redis *redis.Client
}

func NewRedisClient(url string) (*RedisClient, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	return &RedisClient{redis: client}, nil
}

func (r *RedisClient) PublishLog(ctx context.Context, status_code, slug, path, method string) {

	_, err := r.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: "analytics_stream",
		ID:     "*",
		Values: map[string]interface{}{
			"status_code": status_code,
			"path":        path,
			"created_at":  time.Now().Format("2026-29-01 15:04:05"),
			"slug":        slug,
			"method":      method,
		},
	}).Result()

	if err != nil {
		fmt.Println(err)
	}
}
