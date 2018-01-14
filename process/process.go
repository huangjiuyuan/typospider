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

	// For debug only, remove after Elasticsearch finished.
	ContentMap map[string]Content

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
		ContentMap:   make(map[string]Content),

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
				blob := &github.Blob{
					Path: setPath(t.Path, sm.Path),
					Size: *sm.Size,
					SHA:  sm.SHA,
					URL:  sm.URL,
					Data: nil,
				}
				proc.blobqueue.Enqueue(blob)
			} else if sm.Type == "tree" {
				tree, err := proc.Visitor.GetTree(sm.URL)
				tree.Path = setPath(t.Path, sm.Path)
				if err != nil {
					fmt.Printf("[Error] Get tree %s failed\n", t.URL)
				}
				proc.treequeue.Enqueue(tree)
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
	lt := proc.LanguageTool
	cr, err := lt.Check(string(*b.Data), "en", "", "", "", "", "", "", false)
	if err != nil {
		fmt.Println(err)
	}
	if len(cr.Matches) > 0 {
		content, err := NewContent(b.Path, b.Size, b.SHA, b.URL, *b.Data)
		if err != nil {
			fmt.Printf("[Error] Create content %s failed\n", b.Path)
		}
		for _, match := range cr.Matches {
			content.AddTypo(*match)
		}
		proc.ContentMap[content.SHA] = *content
	}
}

func setPath(parent string, current string) string {
	return parent + "/" + current
}
