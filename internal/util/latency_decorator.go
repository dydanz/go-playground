package util

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// GinHandlerLatencyDecorator wraps a Gin handler function and measures its execution time
func GinHandlerLatencyDecorator() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		elapsed := time.Since(start)
		fmt.Printf("Handler %s executed in: %v\n", c.Request.URL.Path, elapsed)
	}
}

// ServiceLatencyDecorator wraps a service method and measures its execution time
func ServiceLatencyDecorator[T any](methodName string, fn func() T) func() T {
	return func() T {
		start := time.Now()
		result := fn()
		elapsed := time.Since(start)
		fmt.Printf("Service method %s executed in: %v\n", methodName, elapsed)
		return result
	}
}

/*
// For Service Methods:
// Use ServiceLatencyDecorator,
func (s *YourService) YourMethod() (ResultType, error) {
    decoratedFn := util.ServiceLatencyDecorator("ServiceName.MethodName", func() ResultType {
        // Your actual logic here
        result, err := s.someOperation()
        if err != nil {
            return nil // or appropriate zero value
        }
        return result
    })

    result := decoratedFn()
    if result == nil {
        return nil, fmt.Errorf("error occurred")
    }
    return result, nil
}

// For Regular Functions:
// Use LatencyDecorator
func YourFunction() {
    decoratedFn := util.LatencyDecorator(func() {
        // Your function logic here
    })
    decoratedFn()
}
*/
