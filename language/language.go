package language

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type LanguageTool struct {
	Addr string
}

type CheckResult struct {
	Software Software `json:"software"`
	Warnings Warnings `json:"warnings"`
	Language Language `json:"language"`
	Matches  []*Match `json:"matches"`
}

type Software struct {
	Name       string  `json:"name"`
	Version    string  `json:"version"`
	BuildDate  string  `json:"buildDate"`
	APIVersion int     `json:"apiVersion"`
	Status     *string `json:"status"`
}

type Warnings struct {
	IncompleteResults bool `json:"incompleteResults"`
}

type Language struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Match struct {
	Message      string         `json:"message"`
	ShortMessage *string        `json:"shortMessage"`
	Offset       int            `json:"offset"`
	Length       int            `json:"length"`
	Replacements []*Replacement `json:"replacements"`
	Context      Context        `json:"context"`
	Sentence     string         `json:"sentence"`
	Rule         Rule           `json:"rule"`
}

type Replacement struct {
	Value *string `json:"value"`
}

type Context struct {
	Text   string `json:"text"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

type Rule struct {
	ID          string   `json:"id"`
	SubID       *string  `json:"subId"`
	Description string   `json:"description"`
	URLs        []*URL   `json:"urls"`
	IssueType   *string  `json:"issueType"`
	Category    Category `json:"category"`
}

type URL struct {
	Value *string `json:"value"`
}

type Category struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
}

type LanguagesResult []Language

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

func (lt *LanguageTool) NewCheckBody(
	text string,
	language string,
	motherTongue string,
	preferredVariants string,
	enabledRules string,
	disabledRules string,
	enabledCategories string,
	disabledCategories string,
	enabledOnly bool) (url.Values, error) {
	if text == "" {
		return nil, fmt.Errorf("missing text parameter")
	}
	if language == "" {
		return nil, fmt.Errorf("missing language parameter")
	}

	strFunc := func(s string) []string {
		if s == "" {
			return nil
		}
		return []string{s}
	}
	boolFunc := func(b bool) []string {
		if b == true {
			return []string{"true"}
		}
		return []string{"false"}
	}
	cb := url.Values{
		"text":               {text},
		"language":           {language},
		"motherTongue":       strFunc(motherTongue),
		"preferredVariants":  strFunc(preferredVariants),
		"enabledRules":       strFunc(enabledRules),
		"disabledRules":      strFunc(disabledRules),
		"enabledCategories":  strFunc(enabledCategories),
		"disabledCategories": strFunc(disabledCategories),
		"enabledOnly":        boolFunc(enabledOnly),
	}

	return cb, nil
}

func (lt *LanguageTool) Check(
	url string,
	text string,
	language string,
	motherTongue string,
	preferredVariants string,
	enabledRules string,
	disabledRules string,
	enabledCategories string,
	disabledCategories string,
	enabledOnly bool) (*CheckResult, error) {
	cb, err := lt.NewCheckBody(
		text,
		language,
		motherTongue,
		preferredVariants,
		enabledRules,
		disabledRules,
		enabledCategories,
		disabledCategories,
		enabledOnly)
	if err != nil {
		return nil, fmt.Errorf("error on creating new check body: %s", err)
	}

	resp, err := http.PostForm(lt.GetURL("", lt.Addr, url), cb)
	if err != nil {
		return nil, fmt.Errorf("error on requesting a check: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error on reading a check response: %s", err)
	}

	cr := new(CheckResult)
	err = json.Unmarshal(body, cr)
	if err != nil {
		return nil, fmt.Errorf("error on parsing a check result: %s", err)
	}

	return cr, nil
}

func (lt *LanguageTool) Languages(url string) (*LanguagesResult, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", lt.GetURL("", lt.Addr, url), nil)
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

	lr := new(LanguagesResult)
	err = json.Unmarshal(body, lr)
	if err != nil {
		return nil, fmt.Errorf("error on parsing a languages result: %s", err)
	}

	return lr, nil
}

func (lt *LanguageTool) GetURL(scheme string, host string, path string) string {
	if scheme == "" {
		scheme = "http"
	}
	if host == "languagetool.org" {
		path = "/api" + path
	}
	return scheme + "://" + host + path
}
