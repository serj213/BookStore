package main

import (
	"fmt"
	"log"

	"github.com/serj213/kafka-service/internal/config"
	"github.com/serj213/kafka-service/internal/kafka"
	"go.uber.org/zap"
)

const (
	dev = "develop"
	local = "local"
	prod = "prod"
)


func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(cfg)

	log, err := initLogger(cfg.Env) 
	if err != nil {
		log.Fatal(err)
	}

	log.Info("logger enabled")
	kafkaCon, err := kafka.NewConsumer(*log,cfg.Kafka.Consumer.BootstapServers, cfg.Kafka.Consumer.GroupId)
	if err != nil {
		log.Fatal(err)
		return
	}
	
	defer kafkaCon.C.Close()

	err = kafkaCon.Subscribe([]string{"newOrder"})
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Debug("previus read message")
	kafkaCon.ReadMessages(cfg.Kafka.Consumer.TimeoutReadMessage)
}

func initLogger(env string) (*zap.SugaredLogger, error) {
	var log *zap.Logger
	var err error

	switch(env){
	case dev:
		log, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}

	case prod:
		log, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}

	default:
		log, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}	
	}

	return log.Sugar(), nil
}