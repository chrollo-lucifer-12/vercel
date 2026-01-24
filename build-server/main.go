package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/supabase-community/supabase-go"
	"golang.org/x/sync/errgroup"
)

type UploadClient struct {
	client *supabase.Client
}

type EnvVars struct {
	API_URL            string
	API_KEY            string
	BUCKET_ID          string
	GIT_REPOSITORY_URL string
}

func NewEnv() (*EnvVars, error) {
	env := &EnvVars{}

	required := []string{
		"GIT_REPOSITORY_URL",
		"API_URL",
		"API_KEY",
		"BUCKET_ID",
	}

	for _, key := range required {
		val, ok := os.LookupEnv(key)
		if !ok || strings.TrimSpace(val) == "" {
			return nil, fmt.Errorf("missing required env var: %s", key)
		}

		switch key {
		case "GIT_REPOSITORY_URL":
			env.GIT_REPOSITORY_URL = val
		case "API_URL":
			env.API_URL = val
		case "API_KEY":
			env.API_KEY = val
		case "BUCKET_ID":
			env.BUCKET_ID = val
		}
	}

	return env, nil
}

func NewUploadClient(apiUrl, apiKey string) (*UploadClient, error) {
	client, err := supabase.NewClient(apiUrl, apiKey, nil)
	if err != nil {
		return nil, err
	}
	return &UploadClient{client: client}, nil
}

func getGitSlug(url string) (string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid git url: %s", url)
	}
	return strings.TrimSuffix(parts[len(parts)-1], ".git"), nil
}

func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func runNpmCommand(dir string, args ...string) error {
	npm := "npm"
	if runtime.GOOS == "windows" {
		npm = "npm.cmd"
	}

	cmd := exec.Command(npm, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (u *UploadClient) uploadFile(ctx context.Context, baseDir, filename, bucketID, slug string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	absFile, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	rel, err := filepath.Rel(baseDir, absFile)
	if err != nil {
		return err
	}

	objectKey := filepath.ToSlash(rel)

	_, err = u.client.Storage.UploadFile(bucketID, slug+"/"+objectKey, file)
	return err
}

func uploadBuild(u *UploadClient, ctx context.Context, bucketID, slug string) error {
	buildPath := filepath.Join("home", "app", "output", "dist")
	if runtime.GOOS == "windows" {
		buildPath = filepath.Join("C:", "home", "app", "output", "dist")
	}

	absBuildDir, err := filepath.Abs(buildPath)
	if err != nil {
		return err
	}

	files := []string{}
	err = filepath.WalkDir(absBuildDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	for _, fp := range files {
		fp := fp

		g.Go(func() error {
			if err := u.uploadFile(ctx, absBuildDir, fp, bucketID, slug); err != nil {
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func main() {
	ctx := context.Background()
	env, err := NewEnv()
	if err != nil {
		fmt.Println("env error:", err)
		return
	}

	slug, err := getGitSlug(env.GIT_REPOSITORY_URL)
	if err != nil {
		fmt.Println("slug error:", err)
		return
	}

	client, err := NewUploadClient(env.API_URL, env.API_KEY)
	if err != nil {
		fmt.Println("supabase client error:", err)
		return
	}

	fmt.Println("Running npm install/build...")

	outputDir := filepath.Join("home", "app", "output")
	if runtime.GOOS == "windows" {
		outputDir = filepath.Join("C:", "home", "app", "output")
	}

	if err := runNpmCommand(outputDir, "install"); err != nil {
		fmt.Println("npm install failed:", err)
		return
	}

	if err := runNpmCommand(outputDir, "run", "build"); err != nil {
		fmt.Println("npm build failed:", err)
		return
	}

	if err := uploadBuild(client, ctx, env.BUCKET_ID, slug); err != nil {
		fmt.Println("upload failed:", err)
		return
	}

	fmt.Println("Upload complete!")
}
