package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string `yaml:"env" env-default:"local"`
	Dsn           string `yaml:"dsn" env-required:"true"`
	MigrationPath string `yaml:"migrationPath" env-required:"true"`
	Grpc          GRPC   `yaml:"grpc" env-required:"true"`
}

type GRPC struct {
	Port    int           `yaml:"port" env-default:"44044"`
	Timeout time.Duration `yaml:"timeout" env-default:"10h"`
}

func GetConfig() (*Config, error) {
	configPath := os.Getenv("configPath")

	if configPath == "" {
		return nil, fmt.Errorf("configPath is empty")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed read config: %w", err)
	}

	return &cfg, nil
}
