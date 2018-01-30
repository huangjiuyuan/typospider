package process

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/huangjiuyuan/typospider/language"
)

type File struct {
	Path      string     `json:"path"`
	Size      int        `json:"size"`
	SHA       string     `json:"sha"`
	URL       string     `json:"url"`
	Fragments []Fragment `json:"fragments"`
	Data      string     `json:"data"`
	Valid     bool       `json:"valid"`
}

type Fragment struct {
	Offset int      `json:"offset"`
	Typos  []string `json:"typos"`
}

type Typo struct {
	SHA    string         `json:"sha"`
	FileID string         `json:"fileId"`
	Match  language.Match `json:"match"`
	Valid  bool           `json:"valid"`
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

func (frag *Fragment) AddTypo(fileId string, match language.Match) (*Typo, error) {
	text := match.Context.Text
	hash := sha1.New()
	hash.Write([]byte(text))
	sha := hex.EncodeToString(hash.Sum(nil))
	frag.Typos = append(frag.Typos, sha)

	typo := &Typo{
		SHA:    sha,
		FileID: fileId,
		Match:  match,
		Valid:  true,
	}

	return typo, nil
}
