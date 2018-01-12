package language

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type LanguageTool struct {
	Addr string
}

type CheckBody struct{}

type CheckResult struct{}

type LanguagesResult struct{}

func NewLanguageTool(host string, port string) (*LanguageTool, error) {
	if host == "" {
		return nil, fmt.Errorf("cannot use an empty host")
	}
	if port == "" {
		return &LanguageTool{
			Addr: host,
		}, nil
	}
	return &LanguageTool{
		Addr: host + ":" + port,
	}, nil
}

func (lt *LanguageTool) NewCheckBody() (*CheckBody, error) {
	return &CheckBody{}, nil
}

func (lt *LanguageTool) Check(url string) (*CheckResult, error) {
	client := &http.Client{}
	cb, err := lt.NewCheckBody()
	if err != nil {
		return nil, fmt.Errorf("error on creating new check body: %s", err)
	}

	data, err := json.Marshal(cb)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", lt.URL("", lt.Addr, url), bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error on creating new request: %s", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error on requesting a check: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error on reading a check response: %s", err)
	}
	fmt.Println(string(body))

	return &CheckResult{}, nil
}

func (lt *LanguageTool) Languages(url string) (*LanguagesResult, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", lt.URL("", lt.Addr, url), nil)
	if err != nil {
		return nil, fmt.Errorf("error on creating new request: %s", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error on requesting languages: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error on reading languages response: %s", err)
	}
	fmt.Println(string(body))

	return &LanguagesResult{}, nil
}

func (lt *LanguageTool) URL(scheme string, host string, path string) string {
	if scheme == "" {
		scheme = "http"
	}
	return scheme + "://" + host + path
}
