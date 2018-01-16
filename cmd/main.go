package main

import (
	"fmt"

	"github.com/huangjiuyuan/typospider/github"
	"github.com/huangjiuyuan/typospider/language"
	"github.com/huangjiuyuan/typospider/process"
)

func main() {
	vis, err := github.NewVisitor(false, "9a0ebc960d050a10f2fbaca29399ee38b12ae61c")
	if err != nil {
		fmt.Println(err)
	}

	lt, err := language.NewLanguageTool("localhost", "6066")
	if err != nil {
		fmt.Println(err)
	}

	proc, err := process.NewProcesser(1000, vis, lt)
	if err != nil {
		fmt.Println(err)
	}

	proc.Elastic.DeleteIndex("kubernetes")
	proc.Elastic.DeleteIndex("typo")

	go proc.ProcessTree("https://api.github.com/repos/kubernetes/kubernetes/git/trees/cec41ac042ea6ac18cf70b7d6f38500b9723e6cb")
	proc.ProcessBlob()
}
