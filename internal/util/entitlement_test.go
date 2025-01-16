package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewCounter tests the creation of a new Counter.
func TestNewCounter(t *testing.T) {
	counter := NewCounter(10)                   // Initialize a new counter with a value of 10
	assert.NotNil(t, counter)                   // Assert that the counter is not nil
	assert.Equal(t, int64(10), counter.Value()) // Assert that the initial value is correct
}

// TestCounter_Take tests the Take method of the Counter.
func TestCounter_Take(t *testing.T) {
	c := NewCounter(10) // Create a new counter with an initial value of 10

	err := c.Take()                      // Attempt to take 1 from the counter
	assert.NoError(t, err)               // Assert that no error occurred
	assert.Equal(t, int64(9), c.Value()) // Assert that the counter value is now 9

	err = c.TakeN(5)                     // Attempt to take 5 from the counter
	assert.NoError(t, err)               // Assert that no error occurred
	assert.Equal(t, int64(4), c.Value()) // Assert that the counter value is now 4
}

// TestCounter_TakeExhausted tests the Take method when the counter is exhausted.
func TestCounter_TakeExhausted(t *testing.T) {
	c := NewCounter(1) // Create a new counter with an initial value of 1

	err := c.Take()        // Take 1 from the counter
	assert.NoError(t, err) // Assert that no error occurred

	err = c.Take()                            // Attempt to take again, which should exhaust the counter
	assert.Equal(t, ErrCounterExhausted, err) // Assert that the correct error is returned
}

// TestCounter_TakeNExhausted tests the TakeN method when the counter is exhausted.
func TestCounter_TakeNExhausted(t *testing.T) {
	c := NewCounter(2) // Create a new counter with an initial value of 2

	err := c.TakeN(2)      // Take 2 from the counter
	assert.NoError(t, err) // Assert that no error occurred

	err = c.TakeN(1)                          // Attempt to take 1 again, which should exhaust the counter
	assert.Equal(t, ErrCounterExhausted, err) // Assert that the correct error is returned
}
