package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-required:"true"`
	Kafka Kafka `yaml:"kafka" env-required:"true"`
}

type Kafka struct {
	Consumer ConsumerKafka `yaml:"consumer" env-required:"true"`
}

type ConsumerKafka struct {
	BootstapServers string `yaml:"bootstrap.servers" env-required:"true"`
	GroupId string `yaml:"group.id" env-required:"true"`
	AutoCommit bool `yaml:"enable.auto.commit" env-default:"true"`
	TimeoutReadMessage time.Duration `yaml:"timeout.read.message" env-default:"5s"`
}



func LoadConfig() (*Config, error) {
	configPath := flag.String("configPath", "./config/local.yaml", "path with config")
	
	if *configPath == "" {
		return nil, fmt.Errorf("config path empty")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}