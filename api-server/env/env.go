package env

import (
	"fmt"
	"os"
	"strings"
)

type EnvVars struct {
	GITHUB_TOKEN   string
	API_URL        string
	API_KEY        string
	REDIS_URL      string
	DSN            string
	AUTH0_DOMAIN   string
	AUTH0_AUDIENCE string
}

func NewEnv() (*EnvVars, error) {
	env := &EnvVars{}

	required := []string{
		"GITHUB_TOKEN",
		"API_URL",
		"API_KEY",
		"DSN",
		"REDIS_URL",
		"AUTH0_DOMAIN",
		"AUTH0_AUDIENCE",
	}

	for _, key := range required {
		val, ok := os.LookupEnv(key)
		if !ok || strings.TrimSpace(val) == "" {
			return nil, fmt.Errorf("missing required env var: %s", key)
		}

		switch key {
		case "GITHUB_TOKEN":
			env.GITHUB_TOKEN = val
		case "API_URL":
			env.API_URL = val
		case "API_KEY":
			env.API_KEY = val
		case "DSN":
			env.DSN = val
		case "REDIS_URL":
			env.REDIS_URL = val
		case "AUTH0_AUDIENCE":
			env.AUTH0_AUDIENCE = val
		case "AUTH0_DOMAIN":
			env.AUTH0_DOMAIN = val

		}
	}

	return env, nil
}
