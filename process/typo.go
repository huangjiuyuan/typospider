package process

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/huangjiuyuan/typospider/language"
)

type File struct {
	Path   string   `json:"path"`
	Size   int      `json:"size"`
	SHA    string   `json:"sha"`
	URL    string   `json:"url"`
	Data   string   `json:"data"`
	Tokens []string `json:"tokens"`
	Typos  []string `json:"typos"`
	Valid  bool     `json:"valid"`
}

type Typo struct {
	SHA   string         `json:"sha"`
	File  string         `json:"file"`
	Match language.Match `json:"match"`
	Valid bool           `json:"valid"`
}

func NewFile(path string, size int, sha string, url string, data []byte) (*File, error) {
	file := new(File)
	file.Path = path
	file.Size = size
	file.SHA = sha
	file.URL = url
	file.Data = string(data)
	file.Valid = true
	return file, nil
}

func (file *File) AddTypo(match language.Match) (*Typo, error) {
	text := match.Context.Text
	hash := sha1.New()
	hash.Write([]byte(text))
	sha := hex.EncodeToString(hash.Sum(nil))
	file.Typos = append(file.Typos, sha)

	typo := &Typo{
		SHA:   sha,
		File:  file.SHA,
		Match: match,
		Valid: true,
	}

	return typo, nil
}
