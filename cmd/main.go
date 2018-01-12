package main

import (
	"fmt"

	"github.com/huangjiuyuan/typospider/api"
)

func main() {
	vis, err := api.NewVisitor(false, 1000, "9b288a1c554007f36d10158705123ba49ee1b8cd")
	if err != nil {
		fmt.Println(err)
	}

	// err = vis.VisitTree("https://api.github.com/repos/kubernetes/kubernetes/git/trees/master")
	go vis.TraverseTree("https://api.github.com/repos/kubernetes/kubernetes/git/trees/cec41ac042ea6ac18cf70b7d6f38500b9723e6cb")

	vis.TraverseBlob()
}
