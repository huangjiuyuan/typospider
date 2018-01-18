package main

import (
	"fmt"

	"github.com/huangjiuyuan/typospider/github"
	"github.com/huangjiuyuan/typospider/language"
	"github.com/huangjiuyuan/typospider/process"
)

func main() {
	vis, err := github.NewVisitor(true, "84521b7efd1c71787383c33ae83ca51387cb27a5")
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

	go proc.ProcessTree("https://api.github.com/repos/kubernetes/kubernetes/git/trees/a071d84d3d2c58c2f704d5c59cd9b254f98f731c")
	proc.ProcessBlob()
}
