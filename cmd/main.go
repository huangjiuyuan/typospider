package main

import (
	"fmt"

	"github.com/huangjiuyuan/typospider/github"
	"github.com/huangjiuyuan/typospider/language"
	"github.com/huangjiuyuan/typospider/process"
)

func main() {
	vis, err := github.NewVisitor(false, "78536a468001488c0f863cf255838071742a792a")
	if err != nil {
		fmt.Println(err)
	}

	lt, err := language.NewLanguageTool("localhost", "6066")
	if err != nil {
		fmt.Println(err)
	}

	cr, err := lt.Check("/v2/check", "I has a mistakes", "en", "", "", "", "", "", "", false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", cr)

	proc, err := process.NewProcesser(1000, vis, lt)
	go proc.ProcessTree("https://api.github.com/repos/kubernetes/kubernetes/git/trees/cec41ac042ea6ac18cf70b7d6f38500b9723e6cb")
	proc.ProcessBlob()
}
