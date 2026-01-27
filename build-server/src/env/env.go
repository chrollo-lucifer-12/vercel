package env

import (
	"fmt"
	"os"
	"strings"
)

type EnvVars struct {
	API_URL            string
	API_KEY            string
	BUCKET_ID          string
	GIT_REPOSITORY_URL string
	SLUG               string
	REDIS_URL          string
	DEPLOYMENT_ID      string
	USER_ENV           string
}

func NewEnv() (*EnvVars, error) {
	env := &EnvVars{}

	required := []string{
		"GIT_REPOSITORY_URL",
		"API_URL",
		"API_KEY",
		"BUCKET_ID",
		"SLUG",
		"REDIS_URL",
		"DEPLOYMENT_ID",
		"USER_ENV",
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
		case "SLUG":
			env.SLUG = val
		case "REDIS_URL":
			env.REDIS_URL = val
		case "DEPLOYMENT_ID":
			env.DEPLOYMENT_ID = val
		case "USER_ENV":
			env.USER_ENV = val
		}
	}

	return env, nil
}
