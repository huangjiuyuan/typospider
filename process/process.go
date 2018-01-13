package process

import (
	"fmt"
	"time"

	"github.com/huangjiuyuan/typospider/github"
	"github.com/huangjiuyuan/typospider/language"
	"github.com/huangjiuyuan/typospider/util/ratelimiter"
)

type Processer struct {
	Visitor      *github.Visitor
	LanguageTool *language.LanguageTool
	Rate         time.Duration

	treequeue ratelimiter.Interface
	blobqueue ratelimiter.Interface
}

func NewProcesser(rate int, vis *github.Visitor, lt *language.LanguageTool) (*Processer, error) {
	if rate < 1000 {
		fmt.Printf("[Warning] API rate exceeded threshold\n")
	}
	p := &Processer{
		Visitor:      vis,
		LanguageTool: lt,
		Rate:         time.Duration(rate) * time.Millisecond,

		treequeue: ratelimiter.New(),
		blobqueue: ratelimiter.New(),
	}
	return p, nil
}

func (proc *Processer) ProcessTree(url string) {
	err := proc.processTree(url)
	if err != nil {
		fmt.Printf("[Error] Processing tree failed: %s\n", err)
	}
}

func (proc *Processer) ProcessBlob() {
	err := proc.processBlob()
	if err != nil {
		fmt.Printf("[Error] Processing blob failed: %s\n", err)
	}
}

func (proc *Processer) processTree(url string) error {
	t, err := proc.Visitor.GetTree(url)
	if err != nil {
		return err
	}
	proc.treequeue.Enqueue(t)

	for {
		item, shutdown := proc.treequeue.Dequeue()
		if shutdown {
			break
		}

		t := item.(*github.Tree)
		for _, sm := range t.Tree {
			if sm.Type == "blob" {
				b := &github.Blob{
					Path: sm.Path,
					SHA:  sm.SHA,
					URL:  sm.URL,
					Data: nil,
				}
				proc.blobqueue.Enqueue(b)
			} else if sm.Type == "tree" {
				t, err := proc.Visitor.GetTree(sm.URL)
				if err != nil {
					fmt.Printf("[Error] Get tree %s failed\n", t.URL)
				}
				proc.treequeue.Enqueue(t)
			}
		}

		if proc.treequeue.Len() == 0 {
			proc.treequeue.ShutDown()
		}
	}

	return nil
}

func (proc *Processer) processBlob() error {
	for {
		item, shutdown := proc.blobqueue.Dequeue()
		if shutdown {
			break
		}

		b := item.(*github.Blob)
		data, err := proc.Visitor.GetBlob(b.URL)
		if err != nil {
			fmt.Printf("[Error] Get blob %s failed\n", b.URL)
		}
		b.Data = &data
		go proc.checkTypo(b)
		if proc.treequeue.ShuttingDown() {
			proc.blobqueue.ShutDown()
		}
	}

	return nil
}

func (proc *Processer) checkTypo(b *github.Blob) {
	fmt.Printf("%#v\n", b)
}
