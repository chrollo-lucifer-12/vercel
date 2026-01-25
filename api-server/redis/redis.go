package redis

import (
	"context"
	"log"

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

func (r *RedisClient) PublishLog(ctx context.Context, log string, channel string) {
	r.redis.Publish(ctx, channel, log)
}

func (r *RedisClient) SubscribeChannel(ctx context.Context, channel string) {
	pubsub := r.redis.Subscribe(ctx, channel)

	defer pubsub.Close()

	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Println("subscribe error:", err)
		return
	}

	ch := pubsub.Channel()

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				log.Println("redis pubsub channel closed")
				return
			}

			log.Printf("received from %s: %s\n", msg.Channel, msg.Payload)

		case <-ctx.Done():
			log.Println("redis subscription context cancelled")
			return
		}
	}
}
