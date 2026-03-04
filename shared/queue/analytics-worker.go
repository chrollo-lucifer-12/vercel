package queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/redis"
	"github.com/hibiken/asynq"
)

type AnalyticsWorker struct {
	server *asynq.Server
	mux    *asynq.ServeMux
	db     *db.DB
	redis  *redis.RedisClient
}

func NewAnalyticsWorker(ctx context.Context, dsn string, redisAddr string) *AnalyticsWorker {
	db, _ := db.NewDB(dsn, ctx)

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: ""},
		asynq.Config{
			Concurrency: 5,
			Queues: map[string]int{
				"analytics": 10,
			},
		},
	)

	mux := asynq.NewServeMux()

	worker := &AnalyticsWorker{
		server: server,
		mux:    mux,
		db:     db,
		redis:  redis.NewRedisClient(redisAddr),
	}

	return worker
}

func (w *AnalyticsWorker) registerHandlers() {
	w.mux.HandleFunc("analytics:track", func(ctx context.Context, t *asynq.Task) error {
		var payload db.WebsiteAnalytics
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}

		data, _ := json.Marshal(payload)
		if err := w.redis.RPush(ctx, "analytics:buffer", data).Err(); err != nil {
			log.Println("Failed to push analytics to Redis buffer:", err)
		}

		return nil
	})

	w.mux.HandleFunc("analytics:flush", func(ctx context.Context, t *asynq.Task) error {
		for {
			events, err := w.redis.LRange(ctx, "analytics:buffer", 0, 9999).Result()
			if err != nil {
				log.Println("Failed to read analytics buffer:", err)
				break
			}
			if len(events) == 0 {
				break
			}

			var batch []db.WebsiteAnalytics
			for _, e := range events {
				var a db.WebsiteAnalytics
				if err := json.Unmarshal([]byte(e), &a); err == nil {
					batch = append(batch, a)
				}
			}

			if len(batch) > 0 {
				if err := w.db.CreateAnalyticsBatch(ctx, batch); err != nil {
					log.Println("Failed to bulk insert analytics:", err)
				}
			}

			if err := w.redis.LTrim(ctx, "analytics:buffer", int64(len(events)), -1).Err(); err != nil {
				log.Println("Failed to trim Redis buffer:", err)
			}
		}
		return nil
	})
}

func (w *AnalyticsWorker) Start() {
	log.Println("Running analytics worker")
	if err := w.server.Run(w.mux); err != nil {
		log.Fatal(err)
	}
}

func (w *AnalyticsWorker) ScheduleDailyFlush(client *asynq.Client) {
	task := asynq.NewTask("analytics:flush", nil)
	_, err := client.Enqueue(task, asynq.MaxRetry(3), asynq.ProcessIn(time.Hour*24))
	if err != nil {
		log.Println("Failed to schedule daily analytics flush:", err)
	}
}
