package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chrollo-lucifer-12/build-server/src/env"
	"github.com/chrollo-lucifer-12/build-server/src/redis"
	"github.com/chrollo-lucifer-12/build-server/src/upload"
	"github.com/chrollo-lucifer-12/build-server/src/utils"
)

func main() {
	ctx := context.Background()
	env, err := env.NewEnv()
	if err != nil {
		fmt.Println("env error:", err)
		return
	}

	userEnv, err := utils.ParseUserEnv(env.USER_ENV)
	if err != nil {
		fmt.Println("user env parse error:", err)
		return
	}

	if err := utils.WriteEnvFile("/home/app/output", userEnv); err != nil {
		fmt.Println("write env file error:", err)
		return
	}

	r, err := redis.NewRedisClient(env.REDIS_URL)
	if err != nil {
		fmt.Println("redis client error:", err)
		return
	}

	client, err := upload.NewUploadClient(env.API_URL, env.API_KEY, r)
	if err != nil {
		fmt.Println("supabase client error:", err)
		return
	}

	fmt.Println("Running npm install/build...")

	outputDir := utils.GetPath([]string{"home", "app", "output"})

	r.PublishLog(ctx, "build started", env.DEPLOYMENT_ID, "INFO")
	if err := utils.RunNpmCommand(
		ctx,
		r,
		env.DEPLOYMENT_ID,
		outputDir,
		"install",
	); err != nil {
		fmt.Println("npm install failed:", err)
		r.PublishLog(ctx, "build failed: "+err.Error(), env.DEPLOYMENT_ID, "ERROR")
		r.PublishLog(ctx, "FAILED", env.DEPLOYMENT_ID, "INFO")
		return
	}
	if err := utils.RunNpmCommand(
		ctx,
		r,
		env.DEPLOYMENT_ID,
		outputDir,
		"run",
		"build",
	); err != nil {
		fmt.Println("npm build failed:", err)
		r.PublishLog(ctx, "build failed: "+err.Error(), env.DEPLOYMENT_ID, "ERROR")
		r.PublishLog(ctx, "FAILED", env.DEPLOYMENT_ID, "INFO")
		return
	}

	if err := client.UploadBuild(ctx, env.BUCKET_ID, env.SLUG); err != nil {
		fmt.Println("upload failed:", err)
		r.PublishLog(ctx, "build failed: "+err.Error(), env.DEPLOYMENT_ID, "ERROR")
		r.PublishLog(ctx, "FAILED", env.DEPLOYMENT_ID, "INFO")
		return
	}

	r.PublishLog(ctx, "upload completed", env.DEPLOYMENT_ID, "INFO")
	fmt.Println("Upload complete!")

	r.PublishLog(ctx, "SUCCESS", env.DEPLOYMENT_ID, "INFO")

	os.Exit(0)
}
