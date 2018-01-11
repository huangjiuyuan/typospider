package api

import (
	"fmt"
	"time"

	"github.com/huangjiuyuan/typospider/util/queue"
)

type Visitor struct {
	Recursive bool
	Rate      time.Duration

	ratelimiter queue.RateLimiter
	queue       queue.WorkQueue
}

type File struct {
	Path string
	Size int32
	SHA  string
	Blob []byte
}

func NewVisitor(recursive bool, rate int) (*Visitor, error) {
	if rate < 1000 {
		fmt.Println("Warning: API rate exceeded threshold.")
	}
	v := &Visitor{
		Recursive: recursive,
		Rate:      time.Duration(rate) * time.Millisecond,
	}
	return v, nil
}

func (v *Visitor) VisitTree(url string, depth int) error {
	t, err := v.GetTree(url)
	if err != nil {
		fmt.Println(err)
	}
	v.ratelimiter.Enqueue(t)
	for v.ratelimiter.Len() > 0 {

	}
	return nil
}
