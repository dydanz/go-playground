package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaService struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

type KafkaConfig struct {
	BrokerURLs []string
}

// NewKafkaService creates a new instance of KafkaService
func NewKafkaService(config KafkaConfig) *KafkaService {
	return &KafkaService{}
}

// PublishMessage publishes a message to a specified topic
func (s *KafkaService) PublishMessage(ctx context.Context, topic string, payload interface{}) error {
	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	// Create writer if not exists
	if s.writer == nil {
		s.writer = &kafka.Writer{
			Addr:     kafka.TCP("localhost:9092"),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}
	}

	// Write message
	err = s.writer.WriteMessages(ctx, kafka.Message{
		Value: jsonData,
	})
	if err != nil {
		return fmt.Errorf("error writing message: %w", err)
	}

	return nil
}

// SubscribeToTopic subscribes to a topic and processes messages
func (s *KafkaService) SubscribeToTopic(ctx context.Context, topic string, groupID string) (<-chan []byte, error) {
	// Create reader if not exists
	if s.reader == nil {
		s.reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{"localhost:9092"},
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		})
	}

	messageChan := make(chan []byte)

	// Start reading messages in a goroutine
	go func() {
		defer close(messageChan)
		defer s.reader.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				message, err := s.reader.ReadMessage(ctx)
				if err != nil {
					fmt.Printf("error reading message: %v\n", err)
					continue
				}
				messageChan <- message.Value
			}
		}
	}()

	return messageChan, nil
}

// Close closes the kafka connections
func (s *KafkaService) Close() error {
	if s.writer != nil {
		if err := s.writer.Close(); err != nil {
			return fmt.Errorf("error closing writer: %w", err)
		}
	}
	if s.reader != nil {
		if err := s.reader.Close(); err != nil {
			return fmt.Errorf("error closing reader: %w", err)
		}
	}
	return nil
}
