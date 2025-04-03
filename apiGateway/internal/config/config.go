package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)


type Config struct {
	Env string `yaml:"env" env-required:"true"`
	GRPC GRPC `yaml:"grpc" env-required:"true"`
	HTTP HTTP `yaml:"http" env-required:"true"`
}

type GRPC struct {
	Port int `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"10h"`
}	

type HTTP struct {
	Addr string `yaml:"addr" env-required:"true"`
}

func Deal() (*Config, error) {

	configPath := flag.String("configPath", "config/local.yaml", "paths configs")

	flag.Parse()

	if *configPath == "" {
		return nil, fmt.Errorf("configPath is empty")
	}
	
	var cfg Config
	
	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed parse config file: %w", err)
	}

	return &cfg, nil

}