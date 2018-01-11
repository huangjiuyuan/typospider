package api

import (
	"fmt"
	"time"

	"github.com/huangjiuyuan/typospider/util/ratelimiter"
)

type Visitor struct {
	Recursive bool
	Rate      time.Duration

	treequeue ratelimiter.Interface
	blobqueue ratelimiter.Interface
}

func NewVisitor(recursive bool, rate int) (*Visitor, error) {
	if rate < 1000 {
		fmt.Println("Warning: API rate exceeded threshold.")
	}
	v := &Visitor{
		Recursive: recursive,
		Rate:      time.Duration(rate) * time.Millisecond,
		treequeue: ratelimiter.New(),
		blobqueue: ratelimiter.New(),
	}
	return v, nil
}

func (vis *Visitor) VisitTree(url string, depth int) error {
	err := vis.visitTree(url, depth)
	if err != nil {
		fmt.Println("Error: Visiting tree failed.")
		return err
	}
	return nil
}

func (vis *Visitor) visitTree(url string, depth int) error {
	t, err := vis.GetTree(url)
	if err != nil {
		return err
	}
	vis.treequeue.Enqueue(t)

	for vis.treequeue.Len() > 0 {
		item := vis.treequeue.Dequeue()
		t := item.(*Tree)
		for _, sm := range t.Tree {
			if sm.Type == "blob" {
				b := &Blob{
					Path: sm.Path,
					SHA:  sm.SHA,
					URL:  sm.URL,
					Data: nil,
				}
				vis.blobqueue.Enqueue(b)
			} else if sm.Type == "tree" {
				t, err := vis.GetTree(sm.URL)
				if err != nil {
					return err
				}
				vis.treequeue.Enqueue(t)
			}
		}
	}

	for vis.blobqueue.Len() > 0 {
		item := vis.blobqueue.Dequeue()
		b := item.(*Blob)
		fmt.Printf("%#v\n", b)
	}

	return nil
}

func (vis *Visitor) visitBlob(url string) error {
	return nil
}
