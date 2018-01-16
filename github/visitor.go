package github

import (
	"net/http"
)

// Visitor is the agent for requesting GitHub API.
type Visitor struct {
	// Whether visiting a tree recursively.
	Recursive bool
	// For authorization.
	Token string
}

// NewVisitor creates a visitor for requesting GitHub API.
func NewVisitor(recursive bool, token string) (*Visitor, error) {
	v := &Visitor{
		Recursive: recursive,
		Token:     token,
	}
	return v, nil
}

// SetAPIAgent sets the request header, including User-Agent, Authorization and Accept fields.
func (vis *Visitor) SetAPIAgent(req *http.Request, raw bool) {
	req.Header.Add("User-Agent", `CCBot`)
	req.Header.Add("Authorization", "token "+vis.Token)
	if raw {
		req.Header.Add("Accept", `application/vnd.github.v3.raw`)
	}
}
