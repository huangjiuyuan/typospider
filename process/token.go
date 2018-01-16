package process

// Tokenizer is for tokenizing raw context.
type Tokenizer struct{}

// NewTokenizer returns a Tokenizer with an error if necessary.
func NewTokenizer() (*Tokenizer, error) {
	return &Tokenizer{}, nil
}

// Tokenize the file context.
func (tk *Tokenizer) Tokenize(file *File) error {
	file.Tokens = append(file.Tokens, file.Data)
	return nil
}
