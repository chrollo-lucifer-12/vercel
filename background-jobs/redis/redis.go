package redis

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	redis  *redis.Client
	ApiURL string
}

type LogRequest struct {
	DeploymentID uuid.UUID      `json:"deployment_id"`
	Log          string         `json:"log"`
	Metadata     map[string]any `json:"metadata"`
	CreatedAt    time.Time      `json:"created_at"`
	Slug         string         `json:"slug"`
}

type UpdateHashRequest struct {
	ProjectID string `json:"project_id"`
	Hash      string `json:"hash"`
	GitURL    string `json:"git_url"`
}

func NewRedisClient(url string, apiURL string) (*RedisClient, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	return &RedisClient{redis: client, ApiURL: apiURL}, nil
}

func (r *RedisClient) ensureGroup(ctx context.Context, stream, group string) {
	err := r.redis.XGroupCreateMkStream(ctx, stream, group, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		log.Fatalf("failed to create group %s: %v", group, err)
	}
}

func (r *RedisClient) sendLogs(ctx context.Context, logs []LogRequest) error {
	if len(logs) == 0 {
		return nil
	}

	data, err := json.Marshal(logs)
	if err != nil {
		log.Println("failed to marshal logs:", err)
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", r.ApiURL+"/api/v1/logs/insert", bytes.NewBuffer(data))
	if err != nil {
		log.Println("failed to create request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("failed to send logs:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("API responded with status:", resp.Status)
	}

	return nil
}

func (r *RedisClient) SubscribeStreams(ctx context.Context, stream string) {

	group := "logs_group"
	consumer := uuid.NewString()

	r.ensureGroup(ctx, stream, group)

	for {
		res, err := r.redis.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{stream, ">"},
			Count:    10,
			Block:    5 * time.Second,
		}).Result()

		if err != nil {
			if err != redis.Nil {
				log.Println("xreadgroup error:", err)
			}
			continue
		}

		var logs []LogRequest
		var msgIDs []string

		for _, s := range res {
			for _, msg := range s.Messages {

				level := msg.Values["level"].(string)
				message := msg.Values["message"].(string)
				createdAt := msg.Values["created_at"].(string)
				deploymentID := msg.Values["deployment_id"].(string)

				t, _ := time.Parse("2006-01-02 15:04:05", createdAt)
				depID, _ := uuid.Parse(deploymentID)

				logs = append(logs, LogRequest{
					DeploymentID: depID,
					Log:          message,
					Metadata:     map[string]any{"level": level, "type": "log"},
					CreatedAt:    t,
				})
				msgIDs = append(msgIDs, msg.ID)

			}
		}

		if len(logs) == 0 {
			continue
		}

		if err := r.sendLogs(ctx, logs); err != nil {
			log.Println("sendLogs failed, skipping ack")
			continue
		}

		r.redis.XAck(ctx, stream, group, msgIDs...)
	}
}

func (r *RedisClient) SubscribeProxyLogs(ctx context.Context, stream string) {
	group := "proxy_group"
	consumer := uuid.NewString()

	r.ensureGroup(ctx, stream, group)

	for {
		res, err := r.redis.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{stream, ">"},
			Count:    10,
			Block:    5 * time.Second,
		}).Result()

		if err != nil {
			if err != redis.Nil {
				log.Println("proxy xreadgroup error:", err)
			}
			continue
		}

		var logs []LogRequest
		var msgIDs []string

		for _, s := range res {
			for _, msg := range s.Messages {

				t, _ := time.Parse(
					"2006-01-02 15:04:05",
					msg.Values["created_at"].(string),
				)

				logs = append(logs, LogRequest{
					Log:       "analytics",
					Slug:      msg.Values["slug"].(string),
					CreatedAt: t,
					Metadata: map[string]any{
						"path":        msg.Values["path"],
						"status_code": msg.Values["status_code"],
						"method":      msg.Values["method"],
					},
				})
				msgIDs = append(msgIDs, msg.ID)

			}
		}

		if len(logs) == 0 {
			continue
		}

		if err := r.sendLogs(ctx, logs); err != nil {
			log.Println("sendLogs failed, skipping ack")
			continue
		}

		r.redis.XAck(ctx, stream, group, msgIDs...)
	}
}

func (r *RedisClient) SubscribeHashStreams(ctx context.Context, stream string) {
	group := "hash_group"
	consumer := uuid.NewString()

	r.ensureGroup(ctx, stream, group)

	for {
		res, err := r.redis.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{stream, ">"},
			Count:    5,
			Block:    10 * time.Second,
		}).Result()

		if err != nil {
			if err != redis.Nil {
				log.Println("hash xreadgroup error:", err)
			}
			continue
		}

		for _, s := range res {
			for _, msg := range s.Messages {

				repoURL := msg.Values["repo_url"].(string)
				lastHash := msg.Values["last_known_hash"].(string)
				projectID := msg.Values["project_id"].(string)

				newHash := GetRepoHash(repoURL)

				if newHash == "" || newHash == lastHash {
					//	r.redis.XAck(ctx, stream, group, msg.ID)
					continue
				}

				body, _ := json.Marshal(UpdateHashRequest{
					ProjectID: projectID,
					Hash:      newHash,
					GitURL:    repoURL,
				})

				resp, err := http.Post(
					r.ApiURL+"/api/v1/hash/update",
					"application/json",
					bytes.NewReader(body),
				)

				if err == nil && resp.StatusCode == http.StatusOK {
					r.redis.XAck(ctx, stream, group, msg.ID)
				}

				if resp != nil {
					resp.Body.Close()
				}
			}
		}
	}
}

func GetRepoHash(repoURL string) string {
	out, err := exec.Command("git", "ls-remote", repoURL, "main").Output()
	if err != nil {
		log.Fatal(err)
		return ""
	}

	var hash string
	parts := strings.Split(string(out), "\t")
	if len(parts) > 0 {
		hash = strings.TrimSpace(parts[0])
	}

	return hash
}
