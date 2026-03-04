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

func (r *RedisClient) RPush(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	return r.client.RPush(ctx, key, value)
}

func (r *RedisClient) LRange(ctx context.Context, key string, start int64, stop int64) *redis.StringSliceCmd {
	return r.client.LRange(ctx, key, start, stop)
}

func (r *RedisClient) LTrim(ctx context.Context, key string, start int64, stop int64) *redis.StatusCmd {
	return r.client.LTrim(ctx, key, start, stop)
}

func (r *RedisClient) DeleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

func (r *RedisClient) Enqueue(ctx context.Context, queue string, payload string) error {
	return r.client.LPush(ctx, queue, payload).Err()
}

func (r *RedisClient) Dequeue(ctx context.Context, queue string) (string, error) {
	result, err := r.client.BRPop(ctx, 0*time.Second, queue).Result()
	if err != nil {
		return "", err
	}
	return result[1], nil
}
