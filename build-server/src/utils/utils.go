package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chrollo-lucifer-12/build-server/src/redis"
)

func GetGitSlug(url string) (string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid git url: %s", url)
	}
	return strings.TrimSuffix(parts[len(parts)-1], ".git"), nil
}

func GetPath(path []string) string {

	dir := filepath.Join(path...)

	if !filepath.IsAbs(dir) {
		dir = string(os.PathSeparator) + dir
	}

	if runtime.GOOS == "windows" {

		if len(path) > 0 && strings.Contains(path[0], ":") {
			return dir
		}

		dir = filepath.Join("C:", dir)
	}

	return dir
}

func RunNpmCommand(
	ctx context.Context,
	redisClient *redis.RedisClient,
	channel string,
	dir string,
	args ...string,
) error {

	npm := "npm"
	if runtime.GOOS == "windows" {
		npm = "npm.cmd"
	}

	cmd := exec.Command(npm, args...)
	cmd.Dir = dir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go publishLogs(ctx, redisClient, channel, "stdout", stdout)

	go publishLogs(ctx, redisClient, channel, "stderr", stderr)

	return cmd.Wait()
}

func publishLogs(
	ctx context.Context,
	redisClient *redis.RedisClient,
	channel string,
	source string,
	reader io.Reader,
) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		msg := fmt.Sprintf("[%s] %s", source, line)

		redisClient.PublishLog(ctx, msg, channel, "INFO")
	}
}
