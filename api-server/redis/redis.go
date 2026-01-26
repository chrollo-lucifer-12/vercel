package redis

import (
	"context"
	"log"
	"time"

	"github.com/chrollo-lucifer-12/api-server/clickhouse"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	redis   *redis.Client
	clickDB *clickhouse.ClickHouseDB
}

func NewRedisClient(url string, clickDB *clickhouse.ClickHouseDB) (*RedisClient, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	return &RedisClient{redis: client, clickDB: clickDB}, nil
}

func (r *RedisClient) SubscribeStreams(ctx context.Context, stream string) {
	lastId := "$"

	for {
		select {
		case <-ctx.Done():
			log.Println("redis stream subscription context cancelled")
			return
		default:
		}

		res, err := r.redis.XRead(ctx, &redis.XReadArgs{
			Streams: []string{stream, lastId},
			Count:   10,
			Block:   0,
		}).Result()

		if err != nil {
			if err == redis.Nil {
				continue
			}
			log.Println("stream read error:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if len(res) == 0 {
			continue
		}

		var logs []clickhouse.Log

		for _, msg := range res[0].Messages {
			lastId = msg.ID
			level, _ := msg.Values["level"]
			message, _ := msg.Values["message"]
			createdAt, _ := msg.Values["created_at"]
			deploymentId, _ := msg.Values["deployment_id"]

			logs = append(logs, clickhouse.Log{
				Level:        level.(string),
				Message:      message.(string),
				CreatedAt:    createdAt.(string),
				DeploymentID: deploymentId.(string),
			})
		}

		if len(logs) > 0 {
			if err := r.clickDB.BatchInsertLogs(ctx, logs); err != nil {
				log.Println("clickhouse insert error:", err)
			}
		}
	}
}
