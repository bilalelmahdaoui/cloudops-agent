package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIAPIKey string
	OpenAIModel  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	model := os.Getenv("OPENAI_MODEL")

	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required")
	}

	if model == "" {
		return nil, fmt.Errorf("OPENAI_MODEL is required")
	}

	return &Config{
		OpenAIAPIKey: apiKey,
		OpenAIModel:  model,
	}, nil
}
