package config

import (
	"errors"
	"os"
)

type Config struct {
	TgToken        string
	DeepSeekAPIKey string
	DataPath       string
}

func New() (*Config, error) {
	token, exists := os.LookupEnv("token")
	if !exists {
		return nil, errors.New("не указан токен для телеграмма")
	}
	deepSeekAPIKey, exists := os.LookupEnv("deepSeekAPIKey")
	if !exists {
		return nil, errors.New("не указан api key")
	}

	dataPath, exists := os.LookupEnv("dataPath")
	if !exists {
		dataPath = "./data.db"
	}

	return &Config{
		TgToken:        token,
		DeepSeekAPIKey: deepSeekAPIKey,
		DataPath:       dataPath,
	}, nil
}
