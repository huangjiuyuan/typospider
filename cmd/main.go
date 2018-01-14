package main

import (
	"fmt"

	"github.com/huangjiuyuan/typospider/github"
	"github.com/huangjiuyuan/typospider/language"
	"github.com/huangjiuyuan/typospider/process"
)

func main() {
	vis, err := github.NewVisitor(false, "7f685a7c723f9a3ff38c0145ebf6d993f7c78d17")
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

	for sha, content := range proc.ContentMap {
		fmt.Printf("SHA: %s, Path: %s, Num: %d\n", sha, content.Path, len(content.Typos))
	}
}
