package file

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"

	kaf "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type FileWriter struct {
	file *os.File
}


func NewFile(path string) (*FileWriter, error){
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err){
			return nil, fmt.Errorf("file is not exist")
		}
	}
	
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error open file: %w", err)
	}

	return &FileWriter{
		file: file,
	}, nil
}

func (f *FileWriter) Write(msg *kaf.Message) error {

	var msgBytes bytes.Buffer

	enc := gob.NewEncoder(&msgBytes)

	if err := enc.Encode(msg); err != nil {
		return fmt.Errorf("failed encode msg: %w", err)
	}

	_, err := f.file.Write(msgBytes.Bytes())
	if err != nil {
		return fmt.Errorf("failed readFile: %w", err)
	}
	return nil
}