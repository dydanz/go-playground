package channel

import (
	"sync"
)

// PubSub is the main structure for our simple pub-sub system.
type PubSub struct {
	// mu is a RWMutex that allows multiple goroutines to read from the subscribers map,
	// but only one goroutine to write to it at any given time.
	mu sync.RWMutex

	// subscribers is a map where the key is a topic (as a string) and the value is
	// a slice of channels. Each channel corresponds to a subscriber listening to that topic.
	subscribers map[string][]chan int
}

// NewPubSub creates a new PubSub instance and initializes its subscribers map.
func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]chan int),
	}
}

// Subscribe allows a subscriber to get updates for a specific topic.
// It returns a channel on which the subscriber will receive these updates.
func (ps *PubSub) Subscribe(topic string) <-chan int {
	// Lock the map for writing.
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Create a new channel for this subscriber.
	ch := make(chan int, 1)

	// Append this subscriber's channel to the slice of channels for the given topic.
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)

	// Return the channel to the subscriber.
	return ch
}

// Publish sends the given value to all subscribers of a specific topic.
func (ps *PubSub) Publish(topic string, value int) {
	// Lock the map for reading.
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	// Iterate over all channels (subscribers) for this topic and send the value.
	for _, subscriber := range ps.subscribers[topic] {
		subscriber <- value
	}
}

// Close removes a specific subscriber from a topic and closes its channel.
func (ps *PubSub) Close(topic string, subCh <-chan int) {
	// Lock the map for writing.
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Find the subscriber's channel in the slice for the given topic.
	subscribers, found := ps.subscribers[topic]
	if !found {
		return
	}

	for i, subscriber := range subscribers {
		if subscriber == subCh {
			// Close the channel.
			close(subscriber)

			// Remove this channel from the slice.
			ps.subscribers[topic] = append(subscribers[:i], subscribers[i+1:]...)
			break
		}
	}
}
