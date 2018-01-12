package main

import (
	"fmt"

	"github.com/huangjiuyuan/typospider/api"
)

func main() {
	vis, err := api.NewVisitor(false, 1000, "a058f59dc639c6ff048a4df0d40a02007077ad70")
	if err != nil {
		fmt.Println(err)
	}

	// err = vis.VisitTree("https://api.github.com/repos/kubernetes/kubernetes/git/trees/master")
	go vis.TraverseTree("https://api.github.com/repos/kubernetes/kubernetes/git/trees/cec41ac042ea6ac18cf70b7d6f38500b9723e6cb")

	vis.TraverseBlob()
}
