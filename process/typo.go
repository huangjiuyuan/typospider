package process

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/huangjiuyuan/typospider/language"
)

type File struct {
	Path  string   `json:"path"`
	Size  int      `json:"size"`
	SHA   string   `json:"sha"`
	URL   string   `json:"url"`
	Data  []byte   `json:"data"`
	Typos []string `json:"typos"`
	Valid bool     `json:"valid"`
}

type Typo struct {
	SHA   string         `json:"sha"`
	Match language.Match `json:"match"`
	Valid bool           `json:"valid"`
}

func NewFile(path string, size int, sha string, url string, data []byte) (*File, error) {
	c := new(File)
	c.Path = path
	c.Size = size
	c.SHA = sha
	c.URL = url
	c.Data = data
	c.Valid = true
	return c, nil
}

func (c *File) AddTypo(match language.Match) (*Typo, error) {
	text := match.Context.Text
	hash := sha1.New()
	hash.Write([]byte(text))
	sha := hex.EncodeToString(hash.Sum(nil))
	c.Typos = append(c.Typos, sha)

	typo := &Typo{
		SHA:   sha,
		Match: match,
		Valid: true,
	}

	return typo, nil
}
