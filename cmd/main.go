package main

import (
	"fmt"

	"github.com/huangjiuyuan/typospider/api"
)

func main() {
	v, err := api.NewVisitor(false, 1000)
	if err != nil {
		fmt.Println(err)
	}
	t, err := v.GetTree("https://api.github.com/repos/kubernetes/kubernetes/git/trees/master")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", t)

	for _, s := range t.Tree {
		if s.Type == "blob" {
			fmt.Printf("%#v\n", s)
			d, err := v.GetBlob(s.URL)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%s\n", d)
		}
	}
}
