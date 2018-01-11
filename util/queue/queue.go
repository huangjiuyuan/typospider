package queue

import (
	"sync"
)

type WorkQueue interface {
	Len() int
	Enqueue(interface{})
	Dequeue() interface{}
}

type workqueue struct {
	Mutex sync.Mutex
}
