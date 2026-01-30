package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestSubscribeHashStreams_LoadCapacity(t *testing.T) {
	// Setup miniredis
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	// Create Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer redisClient.Close()

	// Mock API server with counters
	var apiCalls atomic.Int64
	var successCalls atomic.Int64
	var failedCalls atomic.Int64

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls.Add(1)

		var req UpdateHashRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			failedCalls.Add(1)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		successCalls.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockAPI.Close()

	// Create RedisClient instance
	client := &RedisClient{
		redis:  redisClient,
		ApiURL: mockAPI.URL,
	}

	// Mock getRepoHash to return unique hashes
	oldGetRepoHash := getRepoHash
	getRepoHash = func(repoURL string) string {
		return "new-hash-" + uuid.NewString()[:8]
	}
	defer func() { getRepoHash = oldGetRepoHash }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream := "load_test_stream"

	// TEST CONFIGURATION - Adjust these to test capacity
	totalMessages := 10000 // Total messages to process
	consumerCount := 1     // Number of concurrent consumers

	t.Logf("Starting load test with %d messages and %d consumers", totalMessages, consumerCount)

	// Start timer
	startTime := time.Now()

	// Add messages to stream
	t.Log("Adding messages to stream...")
	addStart := time.Now()
	for i := 0; i < totalMessages; i++ {
		err := redisClient.XAdd(ctx, &redis.XAddArgs{
			Stream: stream,
			Values: map[string]interface{}{
				"repo_url":        fmt.Sprintf("https://github.com/user/repo-%d", i),
				"last_known_hash": fmt.Sprintf("old-hash-%d", i),
				"project_id":      fmt.Sprintf("project-%d", i),
			},
		}).Err()
		if err != nil {
			t.Fatalf("Failed to add message %d: %v", i, err)
		}
	}
	addDuration := time.Since(addStart)
	t.Logf("Added %d messages in %v (%.2f msg/sec)",
		totalMessages, addDuration, float64(totalMessages)/addDuration.Seconds())

	// Start multiple consumers
	t.Logf("Starting %d consumers...", consumerCount)
	for i := 0; i < consumerCount; i++ {
		go client.SubscribeHashStreams(ctx, stream)
	}

	// Monitor progress
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.After(5 * time.Minute) // 5 minute timeout
	lastReported := int64(0)

	for {
		select {
		case <-timeout:
			elapsed := time.Since(startTime)
			current := apiCalls.Load()
			success := successCalls.Load()
			failed := failedCalls.Load()

			t.Logf("\n=== TIMEOUT AFTER %v ===", elapsed)
			t.Logf("Messages processed: %d/%d (%.1f%%)",
				current, totalMessages, float64(current)/float64(totalMessages)*100)
			t.Logf("Success: %d, Failed: %d", success, failed)
			t.Logf("Throughput: %.2f msg/sec", float64(current)/elapsed.Seconds())
			t.Fatalf("Test timed out - only processed %d/%d messages", current, totalMessages)

		case <-ticker.C:
			current := apiCalls.Load()
			if current > lastReported {
				elapsed := time.Since(startTime)
				throughput := float64(current) / elapsed.Seconds()
				progress := float64(current) / float64(totalMessages) * 100

				t.Logf("Progress: %d/%d (%.1f%%) | Throughput: %.2f msg/sec | Elapsed: %v",
					current, totalMessages, progress, throughput, elapsed.Round(time.Second))
				lastReported = current
			}

			// Check if complete
			if current >= int64(totalMessages) {
				elapsed := time.Since(startTime)
				success := successCalls.Load()
				failed := failedCalls.Load()

				t.Logf("\n=== LOAD TEST COMPLETED ===")
				t.Logf("Total messages: %d", totalMessages)
				t.Logf("Consumers: %d", consumerCount)
				t.Logf("Total duration: %v", elapsed)
				t.Logf("API calls: %d", current)
				t.Logf("  - Success: %d", success)
				t.Logf("  - Failed: %d", failed)
				t.Logf("Average throughput: %.2f msg/sec", float64(current)/elapsed.Seconds())
				t.Logf("Messages per consumer: %.2f msg/sec", float64(current)/elapsed.Seconds()/float64(consumerCount))

				// Verify all messages were processed
				if success != int64(totalMessages) {
					t.Errorf("Expected %d successful API calls, got %d", totalMessages, success)
				}

				// Performance expectations (adjust based on your requirements)
				expectedMinThroughput := 50.0 // minimum 50 msg/sec
				actualThroughput := float64(current) / elapsed.Seconds()

				if actualThroughput < expectedMinThroughput {
					t.Logf("WARNING: Throughput %.2f msg/sec is below expected %.2f msg/sec",
						actualThroughput, expectedMinThroughput)
				} else {
					t.Logf("✓ Throughput meets expectations (%.2f >= %.2f msg/sec)",
						actualThroughput, expectedMinThroughput)
				}

				return
			}
		}
	}
}
