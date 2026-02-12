package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func StringToUUID(idStr string) uuid.UUID {
	res, _ := uuid.Parse(idStr)
	return res
}
