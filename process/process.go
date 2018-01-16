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
	Elastic      *Elastic
	Tokenizer    *Tokenizer
	Rate         time.Duration

	treequeue ratelimiter.Interface
	blobqueue ratelimiter.Interface
}

func NewProcesser(rate int, vis *github.Visitor, lt *language.LanguageTool, initialize bool) (*Processer, error) {
	if rate < 1000 {
		fmt.Printf("[Warning] API rate exceeded threshold\n")
	}

	es, err := InitClient("http", "localhost", "9200", initialize)
	if err != nil {
		return nil, err
	}

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

		if t, ok := item.(*github.Tree); ok {
			for _, sm := range t.Tree {
				if sm.Path == "vendor" || sm.Path == "staging" {
					continue
				}

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
						fmt.Printf("[Error] Get tree %s failed: %s\n", t.URL, err)
					}
					proc.treequeue.Enqueue(tree)
				}
			}

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
	err := proc.Elastic.CreateFileIndex("kubernetes")
	if err != nil {
		fmt.Printf("[Error] Create index failed: %s\n", err)
	}

	err = proc.Elastic.CreateTypoIndex("typo")
	if err != nil {
		fmt.Printf("[Error] Create index failed: %s\n", err)
	}

	for {
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
			go proc.processTypo(b)
			if proc.treequeue.ShuttingDown() {
				proc.blobqueue.ShutDown()
			}
		} else {
			fmt.Printf("[Error] Parse blob %#v failed\n", item)
		}
	}

	return nil
}

func (proc *Processer) processTypo(b *github.Blob) {
	file, err := NewFile(b.Path, b.Size, b.SHA, b.URL, *b.Data)
	if err != nil {
		fmt.Printf("[Error] Create file %s failed: %s\n", b.Path, err)
		return
	}

	err = proc.Tokenizer.Tokenize(file)
	if err != nil {
		fmt.Println(err)
	}

	for _, token := range file.Tokens {
		cr, err := proc.LanguageTool.Check(token, "en", "", "", "", "", "", "", false)
		if err != nil {
			fmt.Println(err)
		}

		if len(cr.Matches) > 0 {
			for _, match := range cr.Matches {
				valid := filterTypo(match)
				if valid {
					typo, err := file.AddTypo(*match)
					if err != nil {
						fmt.Printf("[Error] Add typo %s failed: %s\n", match.Context.Text, err)
						continue
					}

					_, err = proc.Elastic.IndexTypo("typo", *typo)
					if err != nil {
						fmt.Printf("[Error] Index typo %s failed: %s\n", typo.Match.Context.Text, err)
						continue
					}
				}
			}
		}
	}

	if len(file.Typos) > 0 {
		_, err = proc.Elastic.IndexFile("kubernetes", *file)
		if err != nil {
			fmt.Printf("[Error] Index file %s failed: %s\n", file.SHA, err)
		}
	}
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
