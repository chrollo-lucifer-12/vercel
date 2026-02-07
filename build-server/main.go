package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chrollo-lucifer-12/shared/env"
	"github.com/chrollo-lucifer-12/shared/upload"
	"github.com/chrollo-lucifer-12/shared/utils"
)

func main() {
	ctx := context.Background()
	// env, err := env.NewEnv()
	// if err != nil {
	// 	fmt.Println("env error:", err)
	// 	return
	// }
	//
	err := env.Load()
	if err != nil {
		panic(err)
	}

	getUserEnv := os.Getenv("USER_ENV")
	userEnv, err := utils.ParseUserEnv(getUserEnv)
	if err != nil {
		panic(err)
	}

	if err := utils.WriteEnvFile("/home/app/output", userEnv); err != nil {
		fmt.Println("write env file error:", err)
		return
	}

	client, err := upload.NewUploadClient(env.SupabaseUrl.GetValue(), env.SupabaseSecret.GetValue())
	if err != nil {
		fmt.Println("supabase client error:", err)
		return
	}

	fmt.Println("Running npm install/build...")

	outputDir := utils.GetPath([]string{"home", "app", "output"})

	//	r.PublishLog(ctx, "build started", env.DEPLOYMENT_ID, "INFO")
	if err := utils.RunNpmCommand(
		ctx,
		outputDir,
		"install",
	); err != nil {
		fmt.Println("npm install failed:", err)
		// r.PublishLog(ctx, "build failed: "+err.Error(), env.DEPLOYMENT_ID, "ERROR")
		// r.PublishLog(ctx, "FAILED", env.DEPLOYMENT_ID, "INFO")
		return
	}
	if err := utils.RunNpmCommand(
		ctx,

		outputDir,
		"run",
		"build",
	); err != nil {
		fmt.Println("npm build failed:", err)
		// r.PublishLog(ctx, "build failed: "+err.Error(), env.DEPLOYMENT_ID, "ERROR")
		// r.PublishLog(ctx, "FAILED", env.DEPLOYMENT_ID, "INFO")
		return
	}

	bucketID := os.Getenv("BUCKET_ID")
	slug := os.Getenv("SLUG")

	if err := client.UploadBuild(ctx, bucketID, slug); err != nil {
		fmt.Println("upload failed:", err)
		// r.PublishLog(ctx, "build failed: "+err.Error(), env.DEPLOYMENT_ID, "ERROR")
		// r.PublishLog(ctx, "FAILED", env.DEPLOYMENT_ID, "INFO")
		return
	}

	//	r.PublishLog(ctx, "upload completed", env.DEPLOYMENT_ID, "INFO")
	fmt.Println("Upload complete!")

	//r.PublishLog(ctx, "SUCCESS", env.DEPLOYMENT_ID, "INFO")

	os.Exit(0)
}
