package main

import (
	"fmt"

	"github.com/huangjiuyuan/typospider/language"
)

func main() {
	// vis, err := api.NewVisitor(false, 1000, "78536a468001488c0f863cf255838071742a792a")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// go vis.TraverseTree("https://api.github.com/repos/kubernetes/kubernetes/git/trees/cec41ac042ea6ac18cf70b7d6f38500b9723e6cb")
	// vis.TraverseBlob()

	// lt, err := language.NewLanguageTool("languagetool.org", "")
	lt, err := language.NewLanguageTool("localhost", "6066")
	if err != nil {
		fmt.Println(err)
	}
	// lt.Languages("/v2/languages")
	lt.Check("/v2/check", "I has a mistakes", "en", "", "", "", "", "", "", false)
}
