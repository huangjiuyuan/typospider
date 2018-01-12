package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/huangjiuyuan/typospider/util/ratelimiter"
)

type Visitor struct {
	Recursive bool
	Rate      time.Duration
	Token     string

	treequeue ratelimiter.Interface
	blobqueue ratelimiter.Interface
}

func NewVisitor(recursive bool, rate int, token string) (*Visitor, error) {
	if rate < 1000 {
		fmt.Printf("[Warning] API rate exceeded threshold\n")
	}
	v := &Visitor{
		Recursive: recursive,
		Rate:      time.Duration(rate) * time.Millisecond,
		Token:     token,
		treequeue: ratelimiter.New(),
		blobqueue: ratelimiter.New(),
	}
	return v, nil
}

func (vis *Visitor) SetAPIAgent(req *http.Request, raw bool) {
	req.Header.Add("User-Agent", `CCBot`)
	req.Header.Add("Authorization", "token "+vis.Token)
	if raw {
		req.Header.Add("Accept", `application/vnd.github.v3.raw`)
	}
}

func (vis *Visitor) TraverseTree(url string) {
	err := vis.processTree(url)
	if err != nil {
		fmt.Printf("[Error] Traversing tree failed: %s\n", err)
	}
}

func (vis *Visitor) TraverseBlob() {
	err := vis.processBlob()
	if err != nil {
		fmt.Printf("[Error] Traversing blob failed: %s\n", err)
	}
}

func (vis *Visitor) processTree(url string) error {
	t, err := vis.GetTree(url)
	if err != nil {
		return err
	}
	vis.treequeue.Enqueue(t)

	for {
		item, shutdown := vis.treequeue.Dequeue()
		if shutdown {
			break
		}

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
					fmt.Printf("[Error] Get tree %s failed\n", t.URL)
				}
				vis.treequeue.Enqueue(t)
			}
		}

		if vis.treequeue.Len() == 0 {
			vis.treequeue.ShutDown()
		}
	}

	return nil
}

func (vis *Visitor) processBlob() error {
	for {
		item, shutdown := vis.blobqueue.Dequeue()
		if shutdown {
			break
		}

		b := item.(*Blob)
		data, err := vis.GetBlob(b.URL)
		if err != nil {
			fmt.Printf("[Error] Get blob %s failed\n", b.URL)
		}
		b.Data = &data
		fmt.Printf("%#v\n", b)
		if vis.treequeue.ShuttingDown() {
			vis.blobqueue.ShutDown()
		}
	}

	return nil
}
