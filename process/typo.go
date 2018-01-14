package process

import (
	"github.com/huangjiuyuan/typospider/language"
)

type Content struct {
	Path  string `json:"path"`
	Size  int    `json:"size"`
	SHA   string `json:"sha"`
	URL   string `json:"url"`
	Data  []byte `json:"data"`
	Typos []Typo `json:"typos"`
	Valid bool   `json:"valid"`
}

type Typo struct {
	Match language.Match
	Valid bool `json:"valid"`
}

func NewContent(path string, size int, sha string, url string, data []byte) (*Content, error) {
	c := new(Content)
	c.Path = path
	c.Size = size
	c.SHA = sha
	c.URL = url
	c.Data = data
	c.Valid = true
	return c, nil
}

func (c *Content) AddTypo(match language.Match) error {
	typo := &Typo{
		Match: match,
		Valid: true,
	}
	c.Typos = append(c.Typos, *typo)
	return nil
}
