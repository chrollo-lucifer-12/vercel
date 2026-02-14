package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chrollo-lucifer-12/build-server/logs"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/upload"
	"github.com/chrollo-lucifer-12/shared/utils"
	"github.com/google/uuid"
)

func StreamLogsToDB(
	ctx context.Context,
	d *db.DB,
	logChan <-chan db.LogEvent,
	done chan<- struct{},
) {

	for {
		select {
		case logEvent, ok := <-logChan:
			if !ok {
				done <- struct{}{}
				return
			}

			err := d.CreateLogEvents(ctx, &[]db.LogEvent{logEvent})
			if err != nil {
				fmt.Println("log insert failed:", err)
			}

		case <-ctx.Done():
			done <- struct{}{}
			return
		}
	}
}

func main() {
	ctx := context.Background()

	dsn := os.Getenv("DSN")
	slug := os.Getenv("SLUG")
	region := os.Getenv("REGION")
	endPoint := os.Getenv("SUPABASE_ENDPOINT")
	supabaseAccessKey := os.Getenv("SUPABASE_ACCESS_KEY")
	supabaseSecret := os.Getenv("SUPABASE_SECRET_KEY")
	getUserEnv := os.Getenv("USER_ENV")
	bucketID := os.Getenv("BUCKET_ID")
	deploymentId := os.Getenv("DEPLOYMENT_ID")
	deploymentIdUUID, _ := uuid.Parse(deploymentId)
	dispatcher := logs.NewLogDispatcher(200)
	done := make(chan struct{})

	userEnv, err := utils.ParseUserEnv(getUserEnv)
	if err != nil {
		panic(err)
	}

	d, err := db.NewDB(dsn, ctx)
	if err != nil {
		panic(err)
	}

	go StreamLogsToDB(ctx, d, dispatcher.Channel(), done)

	if err := utils.WriteEnvFile("/home/app/output", userEnv); err != nil {
		fmt.Println("Write env file error:", err)
		return
	}

	storage, err := upload.NewMinioStorage(endPoint, supabaseAccessKey, region, supabaseSecret, false)
	if err != nil {
		fmt.Println("Supabase client error:", err)
		return
	}

	client := upload.NewUploadClient(storage)

	dispatcher.Push(deploymentIdUUID, "Running npm install/build...")
	outputDir := utils.GetPath([]string{"home", "app", "output"})

	err = RunNpmCommand(ctx, outputDir, dispatcher, deploymentIdUUID, "install")
	if err != nil {
		dispatcher.Push(deploymentIdUUID, "npm install failed: "+err.Error())
		dispatcher.Close()
		<-done
		return
	}

	err = RunNpmCommand(ctx, outputDir, dispatcher, deploymentIdUUID, "run", "build")
	if err != nil {
		dispatcher.Push(deploymentIdUUID, "npm build failed: "+err.Error())
		dispatcher.Close()
		<-done
		return
	}

	if err := client.UploadBuild(ctx, bucketID, slug); err != nil {
		dispatcher.Push(deploymentIdUUID, "build upload failed: "+err.Error())
		dispatcher.Close()
		<-done
		return
	}

	dispatcher.Push(deploymentIdUUID, "Build successful!")
	dispatcher.Close()
	<-done

	os.Exit(0)
}
