package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/chrollo-lucifer-12/api-server/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/datatypes"
)

type RedisClient struct {
	redis *redis.Client
	db    *models.DB
}

func NewRedisClient(url string, db *models.DB) (*RedisClient, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	return &RedisClient{redis: client, db: db}, nil
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

		var logs []models.LogEvent

		for _, msg := range res[0].Messages {
			lastId = msg.ID
			level, _ := msg.Values["level"]
			message, _ := msg.Values["message"]
			createdAt, _ := msg.Values["created_at"]
			deploymentId, _ := msg.Values["deployment_id"]

			createdAtParsed, _ := time.Parse("2006-01-02 15:04:05", createdAt.(string))
			deploymentIdParsed, _ := uuid.Parse(deploymentId.(string))
			metadata, _ := json.Marshal(map[string]any{
				"level": level.(string),
				"type":  "log",
			})

			logs = append(logs, models.LogEvent{
				Log:          message.(string),
				Base:         models.Base{CreatedAt: createdAtParsed},
				DeploymentID: deploymentIdParsed,
				Metadata:     datatypes.JSON(metadata),
			})
		}

		if len(logs) > 0 {
			if err := r.db.CreateLogEvents(ctx, &logs); err != nil {
				log.Println("db logs insert error:", err)
			}
		}
	}
}

func (r *RedisClient) SubscribeProxyLogs(ctx context.Context, stream string) {
	lastId := "$"

	for {
		select {
		case <-ctx.Done():
			log.Println("redis proxy log subscription context cancelled")
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

		var views []models.LogEvent

		for _, msg := range res[0].Messages {
			lastId = msg.ID

			deploymentId, _ := msg.Values["deployment_id"]
			viewDateStr, _ := msg.Values["created_at"]
			statusCode, _ := msg.Values["status_code"]
			path, _ := msg.Values["path"]
			method, _ := msg.Values["method"]

			viewDateParsed, _ := time.Parse("2006-01-02 15:04:05", viewDateStr.(string))
			deploymentIdParsed, _ := uuid.Parse(deploymentId.(string))
			metadata, _ := json.Marshal(map[string]any{
				"path":        path.(string),
				"status_code": statusCode.(string),
				"method":      method.(string),
			})

			views = append(views, models.LogEvent{
				Base:         models.Base{CreatedAt: viewDateParsed},
				DeploymentID: deploymentIdParsed,
				Log:          "analytics",
				Metadata:     datatypes.JSON(metadata),
			})
		}

		if len(views) > 0 {
			if err := r.db.CreateLogEvents(ctx, &views); err != nil {
				log.Println("db logs insert error:", err)
			}
		}
	}
}
