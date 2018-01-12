package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Blob struct {
	Path string  `json:"path"`
	SHA  string  `json:"sha"`
	URL  string  `json:"url"`
	Data *[]byte `json:"data"`
}

func (vis *Visitor) GetBlob(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error on creating new request: %s", err)
	}
	vis.SetAPIAgent(req, true)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error on requesting a blob: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error on reading a blob response: %s", err)
	}
	return body, nil
}
