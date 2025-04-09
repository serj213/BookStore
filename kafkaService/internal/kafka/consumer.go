package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	kaf "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/serj213/kafka-service/internal/file"
	"go.uber.org/zap"
)


type KafkaConsumer struct {
	Log zap.SugaredLogger
	C *kaf.Consumer
	fileWriter *file.FileWriter
}


func NewConsumer(log zap.SugaredLogger,boostrapService, groupId string, file *file.FileWriter) (KafkaConsumer, error){
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
		fileWriter: file,
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

	writeChan := make(chan *kaf.Message, 500)
	errChan := make(chan error, 1)

	fmt.Println("ReadMessages start")

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(){
			defer wg.Done()
			for msg := range writeChan{
				err := c.fileWriter.Write(msg)
				if err != nil {
					errChan <- err
					return
				}
			}

		}()
	}

	batch := make([]*kaf.Message, 0, 100)
	flushBatch := func() error {
		if len(batch) == 0{
			return nil
		}
		select{
		case writeChan <- batch[0]:
			batch = batch[1:]
			return nil
		case <- ctx.Done():
			return ctx.Err()
		}
	}

	for {
		
		select{
		case <-ctx.Done():
			wg.Wait()
			return nil

		default:
			msg, err := c.C.ReadMessage(10 * time.Second)
			if err != nil {
				fmt.Println("лфалф укк, ", err)
				if kafkaErr, ok := err.(kaf.Error); ok {
                    if kafkaErr.Code() == kaf.ErrTimedOut {
						if err := flushBatch(); err != nil {
							return err
						}
                        continue
                    }
                    if kafkaErr.Code() == kaf.ErrAllBrokersDown {
                        time.Sleep(1 * time.Second)
                        continue
                    }
                }
                return fmt.Errorf("kafka read error: %w", err)
			}
			fmt.Println("msg ", msg)
			batch = append(batch, msg)
		}
	}	
}
