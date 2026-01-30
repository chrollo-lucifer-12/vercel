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

func (r *RedisClient) sendLogs(ctx context.Context, logs []LogRequest) {
	if len(logs) == 0 {
		return
	}

	data, err := json.Marshal(logs)
	if err != nil {
		log.Println("failed to marshal logs:", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", r.ApiURL+"/api/v1/logs/insert", bytes.NewBuffer(data))
	if err != nil {
		log.Println("failed to create request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("failed to send logs:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("API responded with status:", resp.Status)
	}
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
			if err != redis.Nil {
				log.Println("stream read error:", err)
				time.Sleep(1 * time.Second)
			}
			continue
		}

		if len(res) == 0 {
			continue
		}

		var logs []LogRequest

		for _, msg := range res[0].Messages {
			lastId = msg.ID
			level, _ := msg.Values["level"]
			message, _ := msg.Values["message"]
			createdAt, _ := msg.Values["created_at"]
			deploymentId, _ := msg.Values["deployment_id"]

			createdAtParsed, _ := time.Parse("2006-01-02 15:04:05", createdAt.(string))

			var depID uuid.UUID
			if deploymentIdStr, ok := deploymentId.(string); ok {
				parsed, err := uuid.Parse(deploymentIdStr)
				if err == nil {
					depID = parsed
				}
			}

			metadata := map[string]any{
				"level": level.(string),
				"type":  "log",
			}

			logs = append(logs, LogRequest{
				DeploymentID: depID,
				Log:          message.(string),
				Metadata:     metadata,
				CreatedAt:    createdAtParsed,
				Slug:         "",
			})
		}

		r.sendLogs(ctx, logs)
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
			if err != redis.Nil {
				log.Println("stream read error:", err)
				time.Sleep(1 * time.Second)
			}
			continue
		}

		if len(res) == 0 {
			continue
		}

		var logs []LogRequest

		for _, msg := range res[0].Messages {
			lastId = msg.ID
			slug, _ := msg.Values["slug"]
			viewDateStr, _ := msg.Values["created_at"]
			statusCode, _ := msg.Values["status_code"]
			path, _ := msg.Values["path"]
			method, _ := msg.Values["method"]

			viewDateParsed, _ := time.Parse("2006-01-02 15:04:05", viewDateStr.(string))

			metadata := map[string]any{
				"path":        path,
				"status_code": statusCode,
				"method":      method,
			}

			var depID uuid.UUID

			logs = append(logs, LogRequest{
				DeploymentID: depID,
				Log:          "analytics",
				Metadata:     metadata,
				CreatedAt:    viewDateParsed,
				Slug:         slug.(string),
			})
		}

		r.sendLogs(ctx, logs)
	}
}

func (r *RedisClient) SubscribeHashStreams(ctx context.Context, stream string) {
	lastId := "$"

	for {
		select {
		case <-ctx.Done():
			log.Println("redis hash log subscription context cancelled")
			return
		default:
		}

		res, err := r.redis.XRead(ctx, &redis.XReadArgs{
			Streams: []string{stream, lastId},
			Count:   10,
			Block:   0,
		}).Result()

		if err != nil {
			if err != redis.Nil {
				log.Println("stream read error:", err)
				time.Sleep(1 * time.Second)
			}
			continue
		}

		if len(res) == 0 {
			continue
		}

		for _, msg := range res[0].Messages {
			lastId = msg.ID
			repoURL := msg.Values["repo_url"]
			lastKnownHash := msg.Values["last_known_hash"]
			projectID := msg.Values["project_id"]

			newHash := GetRepoHash(repoURL.(string))

			if newHash == "" || newHash == lastKnownHash {
				continue
			}

			reqBody := UpdateHashRequest{
				ProjectID: projectID.(string),
				Hash:      newHash,
				GitURL:    repoURL.(string),
			}

			jsonData, err := json.Marshal(reqBody)
			if err != nil {
				log.Println("failed to marshal JSON:", err)
				continue
			}

			resp, err := http.Post(r.ApiURL+"/api/v1/hash/update", "application/json", bytes.NewReader(jsonData))
			if err != nil {
				log.Println("failed to send request to API server:", err)
				continue
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Println("API server responded with status:", resp.Status)
				continue
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
