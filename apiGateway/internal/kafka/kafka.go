package kafka

import (
	"encoding/json"
	"fmt"
	"time"

	k "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type Kafka struct {
	log zap.SugaredLogger
	pr *k.Producer
}

type KafkaMsg struct {
	RequestTime time.Time `json:"request_time"`
	Method string	`json:"method"`
}


func NewKafkaProducer(log zap.SugaredLogger, bootstrapServers string) (*Kafka, error) {
	if bootstrapServers == ""{
		return nil, fmt.Errorf("bootstrapServers empty")
	}

	pr, err := k.NewProducer(&k.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})
	if err != nil {
		return nil, err
	}

	return &Kafka{
		log: log,
		pr: pr,
	}, nil
}

func (kaf Kafka) SendMessage(topic string, msg KafkaMsg) error {

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("kafka msg marhal err: %w", err)
	}

	deliveryChan := make(chan k.Event)

	err = kaf.pr.Produce(&k.Message{
		TopicPartition: k.TopicPartition{
			Topic: &topic,
			Partition: k.PartitionAny,
		},
		Value: msgBytes,
	}, deliveryChan)
	if err != nil {
		return fmt.Errorf("failed send kafka msg: %w", err)
	}

	go func() {
		for e := range deliveryChan {
            m := e.(*k.Message)
            if m.TopicPartition.Error != nil {
                fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
            } else {
                fmt.Printf("Delivered message to %v\n", m.TopicPartition)
            }
        }
	}()

	kaf.pr.Flush(30 * 1000)

	return nil
}