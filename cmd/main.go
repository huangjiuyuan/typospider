package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TreeResp struct {
	SHA       string       `json:"sha"`
	URL       string       `json:"url"`
	Tree      []*Submodule `json:"tree"`
	Truncated bool         `json:"truncated"`
}

type Submodule struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	Size int32  `json:"size"`
	SHA  string `json:"sha"`
	URL  string `json:"url"`
}

func main() {
	resp, err := http.Get("https://api.github.com/repos/kubernetes/kubernetes/git/trees/master")
	if err != nil {
		fmt.Printf("Error on requesting a tree: %s\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error on reading a tree response: %s\n", err)
	}

	t := new(TreeResp)
	err = json.Unmarshal(body, t)
	if err != nil {
		fmt.Printf("Error on parsing a tree: %s\n", err)
	}

	fmt.Printf("%#v\n", t)
	for i, v := range t.Tree {
		fmt.Printf("%d: %#v\n", i, v)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/kubernetes/kubernetes/git/blobs/9974dc685773fd19dcc1d8bfba1f57e6a62b8b3c", nil)
	req.Header.Add("Accept", `application/vnd.github.v3.raw`)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Error on requesting a blob: %s\n", err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error on reading a blob response: %s\n", err)
	}
	fmt.Printf("%s", body)
}
