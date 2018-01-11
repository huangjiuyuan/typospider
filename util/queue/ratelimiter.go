package queue

import (
	"sync"
	"time"
)

type RateLimiter interface {
	Len() int
	When() time.Duration
	// Forget(interface{})
	// NumRequeues(interface{}) int
	Enqueue(interface{})
	Dequeue() interface{}
}

type ratelimiter struct {
	Mutex sync.Mutex
}
