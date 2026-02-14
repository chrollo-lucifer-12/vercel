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
	Env                  EnvKey = "ENV"
	GithubToken          EnvKey = "GITHUB_TOKEN"
	SupabaseUrl          EnvKey = "SUPABASE_URL"
	SupabaseSecret       EnvKey = "SUPABASE_SECRET"
	Dsn                  EnvKey = "DSN"
	SupabaseAccessKey    EnvKey = "SUPABASE_ACCESS_KEY"
	SupabaseAccessSecret EnvKey = "SUPABASE_ACCESS_SECRET"
	Region               EnvKey = "REGION"
	SupabaseEndpoint     EnvKey = "SUPABASE_ENDPOINT"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

func isDevelopmentMode() bool {
	return Env.GetValue() == EnvDevelopment
}

func Load() {
	_ = godotenv.Load(".env")
}
