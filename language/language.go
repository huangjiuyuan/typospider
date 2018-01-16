package language

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// LanguageTool is for visiting languagetool API.
type LanguageTool struct {
	// Addr represents the address of languagetool server.
	Addr string
}

// CheckResult is the response of check request.
type CheckResult struct {
	Software Software `json:"software"`
	Warnings Warnings `json:"warnings"`
	Language Language `json:"language"`
	Matches  []*Match `json:"matches"`
}

// Software is the software information.
type Software struct {
	// Usually "LanguageTool".
	Name string `json:"name"`
	// A version string like "3.3".
	Version string `json:"version"`
	// Date when the software was built.
	BuildDate string `json:"buildDate"`
	// Version of this API response.
	APIVersion int `json:"apiVersion"`
	// An optional warning.
	Status *string `json:"status"`
}

// Warnings containes warning from the server.
type Warnings struct {
	// IncompleteResults represents whether the result is incomplete.
	IncompleteResults bool `json:"incompleteResults"`
}

// Language information.
type Language struct {
	// Language name.
	Name string `json:"name"`
	// ISO 639-1 code like "en", "en-US", or "ca-ES-valencia".
	Code string `json:"code"`
}

// Match represents an error in the text.
type Match struct {
	// Message about the error displayed to the user.
	Message string `json:"message"`
	// An optional shorter version of message.
	ShortMessage *string `json:"shortMessage"`
	// The 0-based character offset of the error in the text.
	Offset int `json:"offset"`
	// The length of the error in characters.
	Length int `json:"length"`
	// Replacements that might correct the error.
	Replacements []*Replacement `json:"replacements"`
	// Context of the error.
	Context Context `json:"context"`
	// The sentence the error occurred in.
	Sentence string `json:"sentence"`
	// Rule violated by the error.
	Rule Rule `json:"rule"`
}

// Replacement that might correct the error.
type Replacement struct {
	// The replacement string.
	Value *string `json:"value"`
}

// Context of the error.
type Context struct {
	// The error and some text to the left and to the right.
	Text string `json:"text"`
	// The 0-based character offset of the error in the text.
	Offset int `json:"offset"`
	// The length of the error in characters in the context.
	Length int `json:"length"`
}

// Rule violated by the error.
type Rule struct {
	// An rule's identifier that's unique for this language.
	ID string `json:"id"`
	// An optional sub identifier of the rule, used when several rules are grouped.
	SubID *string `json:"subId"`
	// Description of the rule.
	Description string `json:"description"`
	// An optional array of URLs with a more detailed description of the error.
	URLs []*URL `json:"urls"`
	// The Localization Quality Issue Type.
	IssueType *string `json:"issueType"`
	// Category represents the error type.
	Category Category `json:"category"`
}

// URL with a more detailed description of the error.
type URL struct {
	// The URL.
	Value *string `json:"value"`
}

// Category represents the error type.
type Category struct {
	// A category's identifier that's unique for this language.
	ID *string `json:"id"`
	// A short description of the category.
	Name *string `json:"name"`
}

// LanguagesResult is the response of languages request.
type LanguagesResult []Language

// NewLanguageTool returns a LanguageTool with an error if necessary.
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

// NewCheckBody returns a body for check request.
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

// Check requests the check API.
func (lt *LanguageTool) Check(
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

	resp, err := http.PostForm(lt.GetURL("", lt.Addr, "/v2/check"), cb)
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

// Languages request the languages API.
func (lt *LanguageTool) Languages() (*LanguagesResult, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", lt.GetURL("", lt.Addr, "/v2/languages"), nil)
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

// GetURL returns the url for sending requests.
func (lt *LanguageTool) GetURL(scheme string, host string, path string) string {
	if scheme == "" {
		scheme = "http"
	}
	if host == "languagetool.org" {
		path = "/api" + path
	}
	return scheme + "://" + host + path
}
