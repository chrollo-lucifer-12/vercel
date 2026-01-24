package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/supabase-community/supabase-go"
)

type UploadClient struct {
	client *supabase.Client
}

func NewUploadClient() *UploadClient {
	client, err := supabase.NewClient(API_URL, API_KEY, nil)
	if err != nil {
		return nil
	}
	return &UploadClient{client: client}
}

func (u *UploadClient) uploadFile(baseDir, filename string) error {

	file, err := os.Open(filename)
	if err != nil {

		return err
	}

	defer file.Close()

	absFile, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	fmt.Println(absFile)

	rel, err := filepath.Rel(baseDir, absFile)
	if err != nil {
		return err
	}

	fmt.Println(rel)
	slug := getGitSlug()
	objectKey := filepath.ToSlash(rel)

	fmt.Println(objectKey)

	_, err = u.client.Storage.UploadFile(BUCKET_ID, slug+"/"+objectKey, file)
	return nil
}

func uploadBuild() {
	u := NewUploadClient()
	buildPath := filepath.Join("home", "app", "output", "dist")
	if runtime.GOOS == "windows" {
		buildPath = filepath.Join("C:", "home", "app", "output", "dist")
	}
	contents, err := os.ReadDir(buildPath)
	if err != nil {
		fmt.Println(err)
	}

	for _, entry := range contents {
		if entry.IsDir() {
			continue
		}
		filePath := filepath.Join(buildPath, entry.Name())

		fmt.Println("Uploading:", entry.Name())

		buildDir := filepath.Join("C:", "home", "app", "output", "dist")
		if runtime.GOOS == "windows" {
			buildDir = filepath.Join("C:", "home", "app", "output", "dist")
		}
		absBuildDir, _ := filepath.Abs(buildDir)
		fmt.Println(absBuildDir)

		if err := u.uploadFile(absBuildDir, filePath); err != nil {
			fmt.Println("upload failed:", err)
		}
	}
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

func getGitSlug() string {
	val, ok := os.LookupEnv("GIT_REPOSITORY_URL")
	if !ok {
		fmt.Println("not found")
	}

	cmds := strings.Split(val, "/")
	projectUrl := strings.Split(cmds[4], ".")[0]
	return projectUrl
}

func main() {

	getGitSlug()

	fmt.Println("executing script.js ....")

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

	uploadBuild()
}
