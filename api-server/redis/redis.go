package redis

import (
	"context"
	"fmt"

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

func (r *RedisClient) PublishLog(ctx context.Context, repoURL, lastKnownHash, projectID string) {

	_, err := r.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: "hash_stream",
		ID:     "*",
		Values: map[string]interface{}{
			"repo_url":        repoURL,
			"last_known_hash": lastKnownHash,
			"project_id":      projectID,
		},
	}).Result()

	if err != nil {
		fmt.Println(err)
	}
}
