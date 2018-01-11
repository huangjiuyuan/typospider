package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Tree struct {
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

func (v *Visitor) GetTree(url string) (*Tree, error) {
	if v.Recursive != true {
		t, err := getTreeUnrecursive(url)
		if err != nil {
			return nil, err
		}
		return t, nil
	}
	t, err := getTreeRecursive(url)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func getTreeUnrecursive(url string) (*Tree, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error on requesting a tree: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error on reading a tree response: %s", err)
	}

	t := new(Tree)
	err = json.Unmarshal(body, t)
	if err != nil {
		return nil, fmt.Errorf("error on parsing a tree: %s", err)
	}

	if t.Truncated != false {
		fmt.Println("Warning: Result of getting url has been truncated.")
	}
	return t, nil
}

func getTreeRecursive(url string) (*Tree, error) {
	resp, err := http.Get(url + "?recursive=1")
	if err != nil {
		return nil, fmt.Errorf("error on requesting a recursive tree: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error on reading a recursive tree response: %s", err)
	}

	t := new(Tree)
	err = json.Unmarshal(body, t)
	if err != nil {
		return nil, fmt.Errorf("error on parsing a recursive tree: %s", err)
	}

	if t.Truncated != false {
		fmt.Println("Warning: Result of getting url has been truncated.")
	}
	return t, nil
}
