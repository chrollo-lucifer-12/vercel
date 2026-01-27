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

func (r *RedisClient) PublishLog(ctx context.Context, log string, deployment_id string, level string) {

	_, err := r.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: "analytics_stream",
		ID:     "*",
		Values: map[string]interface{}{
			"level":         level,
			"message":       log,
			"created_at":    time.Now().UnixMilli(),
			"deployment_id": deployment_id,
		},
	}).Result()

	if err != nil {
		fmt.Println(err)
	}
}
