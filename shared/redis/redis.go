package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(redisURL string) *RedisClient {
	opt, _ := redis.ParseURL(redisURL)
	client := redis.NewClient(opt)
	return &RedisClient{client: client}
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	r.client.Set(ctx, key, value, expiration)

}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {

	fmt.Println("cache for key : ", key)

	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Del(ctx context.Context, key string) {
	r.client.Del(ctx, key)
}
