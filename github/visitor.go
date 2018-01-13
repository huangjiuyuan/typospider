package github

import (
	"net/http"
)

type Visitor struct {
	Recursive bool
	Token     string
}

func NewVisitor(recursive bool, token string) (*Visitor, error) {
	v := &Visitor{
		Recursive: recursive,
		Token:     token,
	}
	return v, nil
}

func (vis *Visitor) SetAPIAgent(req *http.Request, raw bool) {
	req.Header.Add("User-Agent", `CCBot`)
	req.Header.Add("Authorization", "token "+vis.Token)
	if raw {
		req.Header.Add("Accept", `application/vnd.github.v3.raw`)
	}
}
