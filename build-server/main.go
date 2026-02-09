package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/upload"
	"github.com/chrollo-lucifer-12/shared/utils"
	"github.com/google/uuid"
)

func insertLogs(d *db.DB, npmLogs1, npmLogs2 []string, logs []db.LogEvent, slugUUID uuid.UUID, ctx context.Context) {
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		allLogs := append([]db.LogEvent{}, logs...)
		for _, logLine := range npmLogs1 {
			allLogs = append(allLogs, db.LogEvent{DeploymentID: slugUUID, Log: logLine})
		}
		for _, logLine := range npmLogs2 {
			allLogs = append(allLogs, db.LogEvent{DeploymentID: slugUUID, Log: logLine})
		}

		if err := d.CreateLogEvents(ctx, &allLogs); err != nil {
			fmt.Printf("Attempt %d: failed to insert logs: %v\n", attempt, err)
		} else {
			return
		}
	}

}

func main() {
	ctx := context.Background()

	dsn := os.Getenv("DSN")
	slug := os.Getenv("SLUG")
	supabaseUrl := os.Getenv("API_URL")
	supabaseSecret := os.Getenv("API_KEY")
	getUserEnv := os.Getenv("USER_ENV")
	bucketID := os.Getenv("BUCKET_ID")
	deploymentId := os.Getenv("DEPLOYMENT_ID")

	slugUUID, _ := uuid.Parse(deploymentId)

	userEnv, err := utils.ParseUserEnv(getUserEnv)
	if err != nil {
		panic(err)
	}

	var logs []db.LogEvent

	d, err := db.NewDB(dsn, ctx)
	if err != nil {
		panic(err)
	}

	if err := utils.WriteEnvFile("/home/app/output", userEnv); err != nil {
		fmt.Println("Write env file error:", err)
		return
	}

	client, err := upload.NewUploadClient(supabaseUrl, supabaseSecret)
	if err != nil {
		fmt.Println("Supabase client error:", err)
		return
	}

	logs = append(logs, db.LogEvent{DeploymentID: slugUUID, Log: "Running npm install/build..."})
	fmt.Println("Running npm install/build...")

	outputDir := utils.GetPath([]string{"home", "app", "output"})

	npmLogs1, err := utils.RunNpmCommand(ctx, outputDir, "install")
	if err != nil {
		fmt.Println("npm install failed:", err)
		logs = append(logs, db.LogEvent{DeploymentID: slugUUID, Log: "npm install failed: " + err.Error()})
		insertLogs(d, npmLogs1, nil, logs, slugUUID, ctx)
		return
	}

	npmLogs2, err := utils.RunNpmCommand(ctx, outputDir, "run", "build")
	if err != nil {
		fmt.Println("npm build failed:", err)
		logs = append(logs, db.LogEvent{DeploymentID: slugUUID, Log: "npm build failed: " + err.Error()})
		insertLogs(d, npmLogs1, npmLogs2, logs, slugUUID, ctx)
		return
	}

	if err := client.UploadBuild(ctx, bucketID, slug); err != nil {
		fmt.Println("Upload failed:", err)
		logs = append(logs, db.LogEvent{DeploymentID: slugUUID, Log: "Upload failed: " + err.Error()})
		insertLogs(d, npmLogs1, npmLogs2, logs, slugUUID, ctx)
		return
	}

	fmt.Println("Upload complete!")
	logs = append(logs, db.LogEvent{DeploymentID: slugUUID, Log: "Build successful!"})
	insertLogs(d, npmLogs1, npmLogs2, logs, slugUUID, ctx)

	os.Exit(0)
}
