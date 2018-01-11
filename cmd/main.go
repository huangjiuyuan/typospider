package main

import (
	"fmt"

	"github.com/huangjiuyuan/typospider/api"
)

func main() {
	vis, err := api.NewVisitor(false, 1000)
	if err != nil {
		fmt.Println(err)
	}

	// err = vis.VisitTree("https://api.github.com/repos/kubernetes/kubernetes/git/trees/master", 0)
	err = vis.VisitTree("https://api.github.com/repos/kubernetes/kubernetes/git/trees/cec41ac042ea6ac18cf70b7d6f38500b9723e6cb", 0)
	if err != nil {
		fmt.Println(err)
	}
}
