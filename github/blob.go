package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Blob contains metadata of a GitHub blob.
type Blob struct {
	// The GitHub path.
	Path string `json:"path"`
	// Size of a blob.
	Size int `json:"size"`
	// SHA is the identifier.
	SHA string `json:"sha"`
	// URL is for requesting GitHub API.
	URL string `json:"url"`
	// Data contains the raw content of a blob.
	Data *[]byte `json:"data"`
}

// GetBlob gets raw content of a GitHub blob.
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
