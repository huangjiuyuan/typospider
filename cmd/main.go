package main

import (
	"fmt"

	"github.com/huangjiuyuan/typospider/language"
)

func main() {
	// vis, err := api.NewVisitor(false, 1000, "9e44a47bf5ec4f43d2804f0b01c49c50b5e39ed1")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// go vis.TraverseTree("https://api.github.com/repos/kubernetes/kubernetes/git/trees/cec41ac042ea6ac18cf70b7d6f38500b9723e6cb")
	// vis.TraverseBlob()

	lt, err := language.NewLanguageTool("languagetool.org", "")
	if err != nil {
		fmt.Println(err)
	}
	lt.Check("/api/v2/check")
	lt.Languages("/api/v2/languages")
}
