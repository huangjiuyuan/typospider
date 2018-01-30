package main

import (
	"fmt"

	"github.com/huangjiuyuan/typospider/github"
	"github.com/huangjiuyuan/typospider/language"
	"github.com/huangjiuyuan/typospider/process"
)

func main() {
	vis, err := github.NewVisitor(true, "68999f8a97ee7b912fa2b55da098d9a9021c5e04")
	if err != nil {
		fmt.Println(err)
	}

	lt, err := language.NewLanguageTool("localhost", "6066")
	if err != nil {
		fmt.Println(err)
	}

	proc, err := process.NewProcesser(1000, vis, lt, 10, true)
	if err != nil {
		fmt.Println(err)
	}

	go proc.ProcessTree("https://api.github.com/repos/kubernetes/kubernetes/git/trees/a740c006931a59cc99cfbb103208758bbc42baf0")
	proc.ProcessBlob()
}
