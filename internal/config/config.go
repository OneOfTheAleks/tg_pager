package config

import (
	"errors"
	"os"
)

type Config struct {
	TgToken   string
	APIKey    string
	ModelName string
	DataPath  string
	Addr      string
	Port      string
}

func New() (*Config, error) {
	token, exists := os.LookupEnv("token")
	if !exists {
		return nil, errors.New("не указан токен для телеграмма")
	}

	dataPath, exists := os.LookupEnv("dataPath")
	if !exists {
		dataPath = "./data.db"
	}

	apiKey, exists := os.LookupEnv("apiKey")
	if !exists {
		return nil, errors.New("не указан api key")
	}
	modelName, _ := os.LookupEnv("modelName")

	addr, exists := os.LookupEnv("addr")
	if !exists {
		addr = "localhost"
	}

	port, exists := os.LookupEnv("port")
	if !exists {
		port = "8080"
	}

	return &Config{
		TgToken:   token,
		APIKey:    apiKey,
		DataPath:  dataPath,
		ModelName: modelName,
		Addr:      addr,
		Port:      port,
	}, nil
}
