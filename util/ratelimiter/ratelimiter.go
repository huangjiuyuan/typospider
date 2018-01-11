package ratelimiter

import (
	"sync"
)

type Interface interface {
	Len() int
	// When() time.Duration
	// Forget(interface{})
	// NumRequeues(interface{}) int
	Enqueue(interface{})
	Dequeue() interface{}
}

type RateLimiter struct {
	cond     *sync.Cond
	queue    []interface{}
	shutdown bool
}

func New() *RateLimiter {
	return &RateLimiter{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (rl *RateLimiter) Len() int {
	rl.cond.L.Lock()
	defer rl.cond.L.Unlock()
	return len(rl.queue)
}

func (rl *RateLimiter) Enqueue(v interface{}) {
	rl.cond.L.Lock()
	defer rl.cond.L.Unlock()
	rl.queue = append(rl.queue, v)
}

func (rl *RateLimiter) Dequeue() interface{} {
	rl.cond.L.Lock()
	defer rl.cond.L.Unlock()
	var v interface{}
	v, rl.queue = rl.queue[0], rl.queue[1:]
	return v
}
