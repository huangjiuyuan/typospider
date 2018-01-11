package workqueue

import (
	"sync"
)

type Interface interface {
	Len() int
	Enqueue(interface{})
	Dequeue() interface{}
}

type WorkQueue struct {
	cond *sync.Cond
}

func New() *WorkQueue {
	return &WorkQueue{}
}

func (wq *WorkQueue) Len() int {
	// TODO
	return 0
}

func (wq *WorkQueue) Enqueue(interface{}) {
	// TODO
}

func (wq *WorkQueue) Dequeue() interface{} {
	// TODO
	return nil
}
