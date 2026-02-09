package utils

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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
	dir string,
	args ...string,
) ([]string, error) {

	npm := "npm"
	if runtime.GOOS == "windows" {
		npm = "npm.cmd"
	}

	cmd := exec.Command(npm, args...)
	cmd.Dir = dir

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	logs := []string{}

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			logs = append(logs, scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			logs = append(logs, scanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		return logs, err
	}

	return logs, nil
}

func ParseUserEnv(jsonStr string) (map[string]string, error) {
	var envVars map[string]string
	err := json.Unmarshal([]byte(jsonStr), &envVars)
	if err != nil {
		return nil, err
	}
	return envVars, nil
}

func WriteEnvFile(dir string, envVars map[string]string) error {
	path := filepath.Join(dir, ".env")
	content := ""
	for k, v := range envVars {
		content += k + "=" + v + "\n"
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func GetCacheKey(subdomain, path string) string {
	raw := subdomain + ":" + path
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}
