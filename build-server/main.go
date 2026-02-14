package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chrollo-lucifer-12/build-server/logs"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/storage"
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
			fmt.Println(logEvent)
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

	updateDeploymentFunc := func(status string) {
		d.UpdateDeployment(ctx, deploymentIdUUID, db.Deployment{Status: status})
	}

	go StreamLogsToDB(ctx, d, dispatcher.Channel(), done)

	if err := utils.WriteEnvFile("/home/app/output", userEnv); err != nil {
		updateDeploymentFunc("FAILED")
		fmt.Println("Write env file error:", err)
		return
	}

	s, err := storage.NewS3Storage(endPoint, supabaseAccessKey, supabaseSecret, region, bucketID)
	if err != nil {
		updateDeploymentFunc("FAILED")
		fmt.Println(err)
		return
	}

	dispatcher.Push(deploymentIdUUID, "Running npm install/build...")
	outputDir := utils.GetPath([]string{"home", "app", "output"})

	err = RunNpmCommand(ctx, outputDir, dispatcher, deploymentIdUUID, "install")
	if err != nil {
		dispatcher.Push(deploymentIdUUID, "npm install failed: "+err.Error())
		dispatcher.Close()
		<-done
		updateDeploymentFunc("FAILED")
		return
	}

	err = RunNpmCommand(ctx, outputDir, dispatcher, deploymentIdUUID, "run", "build")
	if err != nil {
		dispatcher.Push(deploymentIdUUID, "npm build failed: "+err.Error())
		dispatcher.Close()
		<-done
		updateDeploymentFunc("FAILED")
		return
	}

	if err := s.UploadDirectory(ctx, "/home/app/output/dist", slug, dispatcher, deploymentIdUUID); err != nil {
		fmt.Println("build upload failed: " + err.Error())
		dispatcher.Push(deploymentIdUUID, "build upload failed: "+err.Error())
		dispatcher.Close()
		<-done
		updateDeploymentFunc("FAILED")
		return
	}

	dispatcher.Push(deploymentIdUUID, "Build successful!")
	dispatcher.Close()
	<-done

	updateDeploymentFunc("SUCCESS")

	os.Exit(0)
}
