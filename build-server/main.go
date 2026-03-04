package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/redis"
	"github.com/chrollo-lucifer-12/shared/storage"
	"github.com/chrollo-lucifer-12/shared/utils"
	"github.com/google/uuid"
)

func pushStreamLogsToDB(
	ctx context.Context,
	redisClient *redis.RedisClient,
	d *db.DB,
	streamName string,
	deploymentId uuid.UUID,
) {

	time.Sleep(2 * time.Second)

	messages, err := redisClient.XRange(ctx, streamName, "-", "+").Result()
	if err != nil {
		fmt.Println("Failed to read stream:", err)
		return
	}

	var logs []db.LogEvent

	for _, msg := range messages {
		if m, ok := msg.Values["message"]; ok {
			logs = append(logs, db.LogEvent{
				DeploymentID: deploymentId,
				Log:          m.(string),
			})
		}
	}

	err = d.CreateLogEvents(ctx, &logs)
	if err != nil {
		fmt.Println("Failed to save logs to DB:", err)
		return
	}

	redisClient.Del(ctx, streamName)

	fmt.Println("Logs finalized and stream deleted")
}

func main() {
	ctx := context.Background()

	dsn := os.Getenv("DSN")
	slug := os.Getenv("SLUG")
	region := os.Getenv("REGION")
	endPoint := os.Getenv("SUPABASE_ENDPOINT")
	supabaseAccessKey := os.Getenv("SUPABASE_ACCESS_KEY")
	supabaseSecret := os.Getenv("SUPABASE_SECRET_KEY")
	//	getUserEnv := os.Getenv("USER_ENV")
	bucketID := os.Getenv("BUCKET_ID")
	deploymentId := os.Getenv("DEPLOYMENT_ID")
	redisURL := os.Getenv("REDIS_URL")
	deploymentIdUUID, _ := uuid.Parse(deploymentId)
	streamName := "deployment_logs:" + deploymentId
	// userEnv, err := utils.ParseUserEnv(getUserEnv)
	// if err != nil {
	// 	panic(err)
	// }

	d, err := db.NewDB(dsn, ctx)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewRedisClient(redisURL)

	updateDeploymentFunc := func(status string) {
		d.UpdateDeployment(ctx, deploymentIdUUID, db.Deployment{Status: status})
	}

	logger := func(message string) {
		_, err := redisClient.StreamAdd(ctx, streamName, map[string]interface{}{
			"message": message,
		})
		if err != nil {
			fmt.Println("Failed to write log to Redis:", err)
		}
	}

	s, err := storage.NewS3Storage(endPoint, supabaseAccessKey, supabaseSecret, region, bucketID)
	if err != nil {
		updateDeploymentFunc("FAILED")
		fmt.Println(err)
		return
	}

	logger("Running npm install/build...")
	outputDir := utils.GetPath([]string{"home", "app", "output"})

	err = RunNpmCommand(ctx, outputDir, streamName, logger, "install")
	if err != nil {
		logger("npm install failed: " + err.Error())
		updateDeploymentFunc("FAILED")
		return
	}

	err = RunNpmCommand(ctx, outputDir, streamName, logger, "run", "build")
	if err != nil {
		logger("npm build failed: " + err.Error())
		updateDeploymentFunc("FAILED")
		return
	}

	if err := s.UploadDirectory(ctx, "/home/app/output/dist", slug, deploymentIdUUID, logger); err != nil {
		fmt.Println("build upload failed: " + err.Error())
		logger("build upload failed: " + err.Error())
		updateDeploymentFunc("FAILED")
		return
	}

	logger("build successful!")

	updateDeploymentFunc("SUCCESS")

	os.Exit(0)
}
