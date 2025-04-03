package kafka

import (
	"fmt"
	"time"

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


func (c KafkaConsumer) ReadMessages(timeout time.Duration) error{

	log := c.Log.With(zap.String("struct", "KafkaConsumer"), zap.String("method", "ReadMessages"))

	for {
		msg, err := c.C.ReadMessage(timeout)
		if err != nil {
			return fmt.Errorf("failed kafka readMessage: %v", err)
		}

		log.Debug(
			zap.String("msg-id", string(msg.Key)),
			zap.String("topic", *msg.TopicPartition.Topic),
			zap.Time("timestamp", msg.Timestamp),
		)


		writeChan := make(chan *kaf.Message, 100)

		defer func() {
			close(writeChan)
		}()

		go func ()  {
			for msg := range writeChan{
				err := c.fileWriter.Write(msg)
				if err != nil {
					log.Error("failed write file: %v", err)
				}
			}	
		}()

		writeChan <- msg
	}
}
