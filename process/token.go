package process

import (
	"strings"
	"text/scanner"
)

// Tokenizer is for tokenizing raw text.
type Tokenizer struct{}

// NewTokenizer returns a Tokenizer with an error if necessary.
func NewTokenizer() (*Tokenizer, error) {
	return &Tokenizer{}, nil
}

// Tokenize the text.
func (tokenizer *Tokenizer) Tokenize(text string) (map[int]string, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(text))
	s.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanChars | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments
	tokens := make(map[int]string)

	var line, column int
	var token string
	state, offset := 0, 1

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if tok == scanner.Comment {
			if state == 0 {
				offset = s.Line
				token = s.TokenText()
				state = 1
			} else if state == 1 {
				if s.Position.Line == line+1 && s.Position.Column == column {
					state = 2
					token = token + " " + s.TokenText()
				} else {
					tokens[offset] = token
					offset = s.Line
					token = s.TokenText()
				}
			} else {
				if s.Position.Line == line+1 && s.Position.Column == column {
					token = token + " " + s.TokenText()
				} else {
					state = 1
					tokens[offset] = token
					offset = s.Line
					token = s.TokenText()
				}
			}
			line, column = s.Position.Line, s.Position.Column
		}
	}
	tokens[offset] = token

	return tokens, nil
}
