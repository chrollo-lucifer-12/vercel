package redis

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func TestPublishLog_MultipleMessages(t *testing.T) {

	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client, err := NewRedisClient("redis://" + mr.Addr())
	require.NoError(t, err)

	ctx := context.Background()

	logs := []struct {
		message      string
		deploymentID string
		level        string
	}{
		{"Starting deployment", "deploy-1", "INFO"},
		{"Building application", "deploy-1", "INFO"},
		{"Build failed", "deploy-1", "ERROR"},
		{"Starting deployment", "deploy-2", "INFO"},
		{"Deployment successful", "deploy-2", "SUCCESS"},
	}

	for _, log := range logs {
		client.PublishLog(ctx, log.message, log.deploymentID, log.level)
	}

	messages, err := client.redis.XRange(ctx, "logs_stream", "-", "+").Result()
	require.NoError(t, err)
	assert.Equal(t, len(logs), len(messages))

	for i, msg := range messages {
		assert.Equal(t, logs[i].level, msg.Values["level"])
		assert.Equal(t, logs[i].message, msg.Values["message"])
		assert.Equal(t, logs[i].deploymentID, msg.Values["deployment_id"])
	}
}
