package github

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
	Size *int   `json:"size"`
	SHA  string `json:"sha"`
	URL  string `json:"url"`
}

func (vis *Visitor) GetTree(url string) (*Tree, error) {
	if vis.Recursive != true {
		t, err := vis.getTreeUnrecursive(url)
		if err != nil {
			return nil, err
		}
		return t, nil
	}
	t, err := vis.getTreeRecursive(url)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (vis *Visitor) getTreeUnrecursive(url string) (*Tree, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error on creating new request: %s", err)
	}
	vis.SetAPIAgent(req, false)
	resp, err := client.Do(req)
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
		fmt.Printf("[Warning] Result has been truncated\n")
	}
	return t, nil
}

func (vis *Visitor) getTreeRecursive(url string) (*Tree, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url+"?recursive=1", nil)
	if err != nil {
		return nil, fmt.Errorf("error on creating new request: %s", err)
	}
	vis.SetAPIAgent(req, false)
	resp, err := client.Do(req)
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
		fmt.Printf("[Warning] Result has been truncated\n")
	}
	return t, nil
}