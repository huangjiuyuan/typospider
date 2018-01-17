package process

import (
	"fmt"
	"sync"
	"time"

	"github.com/huangjiuyuan/typospider/github"
	"github.com/huangjiuyuan/typospider/language"
	"github.com/huangjiuyuan/typospider/util/ratelimiter"
)

// Processer contains a Visitor to visit GitHub API, a LanguageTool to check the texts, a Elastic
// agent to operate on a Elasticsearch server, and a Tokenizer to tokenize context of a file. It produces
// trees by visiting GitHub API, and produces blobs by consuming trees it produces.
type Processer struct {
	// Visitor to visit GitHub API.
	Visitor *github.Visitor
	// LanguageTool to check the texts.
	LanguageTool *language.LanguageTool
	// Elastic agent to operate on a Elasticsearch server.
	Elastic *Elastic
	// Tokenizer to tokenize context of a file.
	Tokenizer *Tokenizer
	// Rate of the GitHub visitor.
	Rate time.Duration

	// Wait for goroutines to finish.
	wg sync.WaitGroup
	// Keep concurrency under control
	sema chan struct{}
	// Thread safe rate limiting queue for processing trees.
	treequeue ratelimiter.Interface
	// Thread safe rate limiting queue for processing blobs.
	blobqueue ratelimiter.Interface
}

// NewProcesser returns a Processer with an error if necessary.
func NewProcesser(rate int, vis *github.Visitor, lt *language.LanguageTool, concurrency int, initialize bool) (*Processer, error) {
	if rate < 1000 {
		fmt.Printf("[Warning] API rate exceeded threshold\n")
	}

	// Create the Elasticsearch client.
	es, err := InitClient("http", "localhost", "9200", initialize)
	if err != nil {
		return nil, err
	}

	// Create the tokenizer.
	tk, err := NewTokenizer()
	if err != nil {
		return nil, err
	}

	p := &Processer{
		Visitor:      vis,
		LanguageTool: lt,
		Elastic:      es,
		Tokenizer:    tk,
		Rate:         time.Duration(rate) * time.Millisecond,

		wg:        sync.WaitGroup{},
		sema:      make(chan struct{}, concurrency),
		treequeue: ratelimiter.New(),
		blobqueue: ratelimiter.New(),
	}

	return p, nil
}

// ProcessTree wraps the tree processing function.
func (proc *Processer) ProcessTree(url string) {
	err := proc.processTree(url)
	if err != nil {
		fmt.Printf("[Error] Processing tree failed: %s\n", err)
	}
}

// ProcessBlob wraps the blob processing function.
func (proc *Processer) ProcessBlob() {
	err := proc.processBlob()
	if err != nil {
		fmt.Printf("[Error] Processing blob failed: %s\n", err)
	}
}

func (proc *Processer) processTree(url string) error {
	// Produce a tree then enqueue to the tree queue.
	t, err := proc.Visitor.GetTree(url)
	if err != nil {
		return err
	}
	proc.treequeue.Enqueue(t)

	for {
		// Shut down if received a signal from dequeue operation.
		item, shutdown := proc.treequeue.Dequeue()
		if shutdown {
			break
		}

		if t, ok := item.(*github.Tree); ok {
			for _, sm := range t.Tree {
				// Skip "vendor" and "staging" folders.
				if sm.Path == "vendor" || sm.Path == "staging" {
					continue
				}

				if sm.Type == "blob" {
					// Produce a blob then enqueue to the blob queue if the submodule is a blob.
					blob := &github.Blob{
						Path: setPath(t.Path, sm.Path),
						Size: *sm.Size,
						SHA:  sm.SHA,
						URL:  sm.URL,
						Data: nil,
					}
					proc.blobqueue.Enqueue(blob)
				} else if sm.Type == "tree" {
					// Produce a tree then enqueue to the tree queue if the submodule is a tree.
					tree, err := proc.Visitor.GetTree(sm.URL)
					tree.Path = setPath(t.Path, sm.Path)
					if err != nil {
						fmt.Printf("[Error] Get tree %s failed: %s\n", t.URL, err)
					}
					proc.treequeue.Enqueue(tree)
				}
			}

			// Send a signal if the tree queue is done and no tree is produced.
			if proc.treequeue.Len() == 0 {
				proc.treequeue.ShutDown()
			}
		} else {
			fmt.Printf("[Error] Parse tree %#v failed\n", item)
		}
	}

	return nil
}

func (proc *Processer) processBlob() error {
	// Create the project index.
	err := proc.Elastic.CreateFileIndex("kubernetes")
	if err != nil {
		fmt.Printf("[Error] Create index failed: %s\n", err)
	}

	// Create the typo index.
	err = proc.Elastic.CreateTypoIndex("typo")
	if err != nil {
		fmt.Printf("[Error] Create index failed: %s\n", err)
	}

	for {
		// Shut down if received a signal from dequeue operation.
		item, shutdown := proc.blobqueue.Dequeue()
		if shutdown {
			break
		}

		if b, ok := item.(*github.Blob); ok {
			data, err := proc.Visitor.GetBlob(b.URL)
			if err != nil {
				fmt.Printf("[Error] Get blob %s failed: %s\n", b.URL, err)
			}
			b.Data = &data

			// Block until the semaphore has room. If the concurrency is under control, process the
			// typo produced by the blob.
			proc.wg.Add(1)
			proc.sema <- struct{}{}
			go proc.processTypo(b)

			// Send a signal if the tree queue is done and no tree is produced.
			if proc.treequeue.ShuttingDown() {
				proc.blobqueue.ShutDown()
			}
		} else {
			fmt.Printf("[Error] Parse blob %#v failed\n", item)
		}
	}
	proc.wg.Wait()
	close(proc.sema)

	return nil
}

func (proc *Processer) processTypo(b *github.Blob) {
	// Create a file from the blob.
	file, err := NewFile(b.Path, b.Size, b.SHA, b.URL, *b.Data)
	if err != nil {
		fmt.Printf("[Error] Create file %s failed: %s\n", b.Path, err)
		return
	}

	// Tokenize the file context.
	tokens, err := proc.Tokenizer.Tokenize(file)
	if err != nil {
		fmt.Println(err)
	}

	// Check error for each token.
	for offset, token := range tokens {
		cr, err := proc.LanguageTool.Check(token, "en", "", "", "", "", "", "", false)
		if err != nil {
			fmt.Println(err)
		}

		if len(cr.Matches) > 0 {
			frag := Fragment{offset, []string{}}
			for _, match := range cr.Matches {
				// Filter out any invalid typo.
				valid := filterTypo(match)
				if valid {
					// Add a typo to the file fragment if it is valid.
					typo, err := frag.AddTypo(file.SHA, *match)
					if err != nil {
						fmt.Printf("[Error] Add typo %s failed: %s\n", match.Context.Text, err)
						continue
					}

					// Index the typo to Elasticsearch.
					_, err = proc.Elastic.IndexTypo("typo", *typo)
					if err != nil {
						fmt.Printf("[Error] Index typo %s failed: %s\n", typo.Match.Context.Text, err)
						continue
					}
				}
			}
			if len(frag.Typos) > 0 {
				file.Fragments = append(file.Fragments, frag)
			}
		}
	}

	// If the file contains any fragment, index the file to Elasticsearch.
	if len(file.Fragments) > 0 {
		_, err = proc.Elastic.IndexFile("kubernetes", *file)
		if err != nil {
			fmt.Printf("[Error] Index file %s failed: %s\n", file.SHA, err)
		}
	}

	proc.wg.Done()
	<-proc.sema
}

func filterTypo(match *language.Match) bool {
	if match.Rule.ID == "EN_QUOTES" {
		return false
	} else if match.Rule.ID == "SENTENCE_WHITESPACE" {
		return false
	} else if match.Rule.ID == "COMMA_PARENTHESIS_WHITESPACE" {
		return false
	} else if match.Rule.ID == "UPPERCASE_SENTENCE_START" {
		return false
	} else if match.Rule.ID == "WHITESPACE_RULE" {
		return false
	} else if match.Rule.ID == "ENGLISH_WORD_REPEAT_BEGINNING_RULE" {
		return false
	} else if match.Rule.ID == "DASH_RULE" {
		return false
	} else if match.Rule.ID == "DOUBLE_PUNCTUATION" {
		return false
	} else if match.Rule.ID == "SENTENCE_FRAGMENT" {
		return false
	} else if match.Rule.ID == "EN_UNPAIRED_BRACKETS" {
		return false
	} else if match.Rule.ID == "ENGLISH_WORD_REPEAT_RULE" {
		return false
	} else if match.Rule.ID == "THE_SENT_END" {
		return false
	}

	return true
}

func setPath(parent string, current string) string {
	return parent + "/" + current
}
