package ratelimiter

import (
	"sync"
)

type Interface interface {
	Len() int
	// When() time.Duration
	// Forget(item interface{})
	// NumRequeues(item interface{}) int
	Enqueue(item interface{})
	Dequeue() (item interface{}, shutdown bool)
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

func (rl *RateLimiter) Enqueue(item interface{}) {
	rl.cond.L.Lock()
	defer rl.cond.L.Unlock()
	rl.queue = append(rl.queue, item)
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

	var item interface{}
	item, rl.queue = rl.queue[0], rl.queue[1:]
	return item, false
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
