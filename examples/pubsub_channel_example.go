import (
	"fmt"
	"sync"

	"go-playground/pkg/channel"
)

// Main function demonstrating the pub-sub functionality.
func main() {
	// Create a new pub-sub instance.
	ps := NewPubSub()

	// A subscriber subscribes to "topic1".
	subscriber := ps.Subscribe("topic1")

	// A publisher publishes the value 42 to "topic1".
	go func() {
		ps.Publish("topic1", 42)
	}()

	// The subscriber receives the value.
	value := <-subscriber
	fmt.Println("Received value:", value) // Expected: Received value: 42

	// Close the subscriber's channel and remove it from the topic's subscribers.
	ps.Close("topic1", subscriber)
}
