package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/serj213/kafka-service/internal/config"
	"github.com/serj213/kafka-service/internal/file"
	"github.com/serj213/kafka-service/internal/kafka"
	"go.uber.org/zap"
)

const (
	dev = "develop"
	local = "local"
	prod = "prod"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
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

	file, err := file.NewFile(filepath.Join("../", "logs.txt"))
	if err != nil {
		log.Fatal(err)
	}

	kafkaCon, err := kafka.NewConsumer(*log,cfg.Kafka.Consumer.BootstapServers, cfg.Kafka.Consumer.GroupId, file)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = kafkaCon.Subscribe([]string{"books"})
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Debug("previus read message")

	err = kafkaCon.ReadMessages(ctx)


	stopped := make(chan struct{})
	go func(){
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		defer cancel()
		kafkaCon.C.Close()
		close(stopped)
	}()

	
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