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

func (vi *Visitor) GetBlob(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", `application/vnd.github.v3.raw`)
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
