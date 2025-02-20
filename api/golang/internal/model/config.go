package model

import (
	"fmt"
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	SolcastEndpoint string `env:"SOLCAST_ENDPOINT,required"`
	SolcastApiKey   string `env:"SOLCAST_API_KEY,required"`

	FirebaseProjectId       string `env:"FIREBASE_PROJECT_ID,required"`
	FirebaseDatabaseUrl     string `env:"FIREBASE_DATABASE_URL,required"`
	FirebaseCredentialsPath string `env:"FIREBASE_CREDENTIALS_FILE_PATH,required"`

	SmaBaseUrl string `env:"SMA_BASE_URL"`
}

func LoadEnvs(cfg *Config) error {
	if path := os.Getenv("ENV_FILE_PATH"); path != "" {
		err := godotenv.Load(path)
		if err != nil {
			return fmt.Errorf("error loading .env file: %w", err)
		}
	}
	err := env.Parse(cfg)
	if err != nil {
		return fmt.Errorf("failed to load envs: %w", err)
	}
	return nil
}
