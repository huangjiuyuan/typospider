package main

import (
	"fmt"

	"github.com/huangjiuyuan/typospider/github"
	"github.com/huangjiuyuan/typospider/language"
	"github.com/huangjiuyuan/typospider/process"
)

func main() {
	vis, err := github.NewVisitor(false, "351a6bbaa7a0c9273aa6c3b3666eb73c677b10cf")
	if err != nil {
		fmt.Println(err)
	}

	lt, err := language.NewLanguageTool("localhost", "6066")
	if err != nil {
		fmt.Println(err)
	}

	proc, err := process.NewProcesser(1000, vis, lt)
	go proc.ProcessTree("https://api.github.com/repos/kubernetes/kubernetes/git/trees/cec41ac042ea6ac18cf70b7d6f38500b9723e6cb")
	proc.ProcessBlob()

	for url, content := range proc.ContentMap {
		fmt.Printf("URL: %s, Num of typos: %d\n", url, len(content.Typos))
	}
}
