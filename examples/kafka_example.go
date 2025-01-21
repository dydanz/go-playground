package kafka_example

import (
	"context"
	"fmt"
	"go-playground/pkg/kafka"
	"log"
	"time"
)

func main() {
	// Initialize Kafka service
	kafkaService := kafka.NewKafkaService(kafka.KafkaConfig{
		BrokerURLs: []string{"localhost:9092"},
	})
	defer kafkaService.Close()

	ctx := context.Background()
	topic := "test-topic"

	// Example publisher
	go func() {
		for i := 0; i < 5; i++ {
			payload := map[string]interface{}{
				"message":   fmt.Sprintf("Test message %d", i),
				"timestamp": time.Now(),
			}

			err := kafkaService.PublishMessage(ctx, topic, payload)
			if err != nil {
				log.Printf("Error publishing message: %v", err)
			}
			time.Sleep(time.Second)
		}
	}()

	// Example subscriber
	messageChan, err := kafkaService.SubscribeToTopic(ctx, topic, "test-group")
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	// Process messages
	for message := range messageChan {
		fmt.Printf("Received message: %s\n", string(message))
	}
}
