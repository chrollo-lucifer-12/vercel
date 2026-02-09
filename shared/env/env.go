package env

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvKey string

func (key EnvKey) GetValue() string {
	return os.Getenv(string(key))
}

const (
	Env            EnvKey = "ENV"
	GithubToken    EnvKey = "GITHUB_TOKEN"
	SupabaseUrl    EnvKey = "SUPABASE_URL"
	SupabaseSecret EnvKey = "SUPABASE_SECRET"
	Dsn            EnvKey = "DSN"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

func isDevelopmentMode() bool {
	return Env.GetValue() == EnvDevelopment
}

func Load() error {
	_, err := os.Lstat(".env")
	if err != nil {
		return err
	}

	if err := godotenv.Load(".env"); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
