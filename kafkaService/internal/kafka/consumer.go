package kafka

import (
	"context"
	"fmt"

	kaf "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/serj213/kafka-service/internal/file"
	"go.uber.org/zap"
)


type KafkaConsumer struct {
	Log zap.SugaredLogger
	C *kaf.Consumer
	fileWriter file.FileWriter
}


func NewConsumer(log zap.SugaredLogger,boostrapService, groupId string) (KafkaConsumer, error){
	c, err := kaf.NewConsumer(&kaf.ConfigMap{
		"bootstrap.servers": boostrapService,
		"group.id": groupId,
	})
	if err != nil {
		return KafkaConsumer{}, err
	}

	return KafkaConsumer{
		C: c,
		Log: log,
	}, nil
}



func (c KafkaConsumer) Subscribe(topics []string) error{

	err := c.C.SubscribeTopics(topics, nil)
	if err != nil {
		return fmt.Errorf("failed subscribe topics: %v", err)
	}

	return nil
}


func (c KafkaConsumer) ReadMessages(ctx context.Context) error{

	log := c.Log.With(zap.String("struct", "KafkaConsumer"), zap.String("method", "ReadMessages"))
	writeChan := make(chan *kaf.Message, 100)

	defer close(writeChan)
	
	go func ()  {
		for msg := range writeChan{
			err := c.fileWriter.Write(msg)
			if err != nil {
				log.Error("failed write file: %v", err)
			}
		}	
	}()

	for {
		select{
		case <- ctx.Done():
			log.Info("cancel kafka consumer")
			return nil
		default:
			msg, err := c.C.ReadMessage(-1)
			if err != nil {
				if kafkaErr, ok := err.(kaf.Error); ok && kafkaErr.Code() == kaf.ErrTimedOut {
					continue
				}
				return fmt.Errorf("failed kafka readMessage: %v", err)
			}

			log.Debug(
				zap.String("msg-id", string(msg.Key)),
				zap.String("topic", *msg.TopicPartition.Topic),
				zap.Time("timestamp", msg.Timestamp),
			)

			select {
            case writeChan <- msg:
            case <-ctx.Done():
                return nil
            }
		}	
	}
}
