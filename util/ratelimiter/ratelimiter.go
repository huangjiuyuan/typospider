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
	Dequeue() (interface{}, bool)
	ShutDown()
	ShuttingDown() bool
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
	rl.cond.Signal()
}

func (rl *RateLimiter) Dequeue() (interface{}, bool) {
	rl.cond.L.Lock()
	defer rl.cond.L.Unlock()
	for len(rl.queue) == 0 && !rl.shutdown {
		rl.cond.Wait()
	}
	if len(rl.queue) == 0 {
		return nil, true
	}

	var v interface{}
	v, rl.queue = rl.queue[0], rl.queue[1:]
	return v, false
}

func (rl *RateLimiter) ShutDown() {
	rl.cond.L.Lock()
	defer rl.cond.L.Unlock()
	rl.shutdown = true
	rl.cond.Broadcast()
}

func (rl *RateLimiter) ShuttingDown() bool {
	rl.cond.L.Lock()
	defer rl.cond.L.Unlock()
	return rl.shutdown
}
