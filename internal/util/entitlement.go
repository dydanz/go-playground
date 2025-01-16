package util

import (
	"errors"
)

// Entitlement represents the entitlement structure with teams information.
type Entitlement struct {
	Teams struct {
		Enabled bool    `json:"enabled,omitempty"` // Indicates if teams are enabled
		Size    Counter `json:"size,omitempty"`    // Represents the size of the counter for teams
	} `json:"teams,omitempty"`
}

// ErrCounterExhausted is an error returned when the counter is exhausted.
var ErrCounterExhausted = errors.New("counter exhausted")

// Counter is a type that represents a counter for managing resources.
type Counter int64

// NewCounter initializes a new Counter with the given value.
func NewCounter(i int64) *Counter {
	c := Counter(i)
	return &c
}

// Value returns the current value of the counter.
func (c *Counter) Value() int64 {
	if *c == 0 {
		return 0
	}
	return int64(*c)
}

// Take decrements the counter by 1 and returns an error if exhausted.
func (c *Counter) Take() error {
	return c.TakeN(1) // Calls TakeN with a value of 1
}

// TakeN decrements the counter by the specified amount and returns an error if exhausted.
func (c *Counter) TakeN(i int64) error {
	if *c <= 0 {
		return ErrCounterExhausted // Return error if counter is exhausted
	}

	*c -= Counter(i) // Decrement the counter by the specified amount
	return nil
}
