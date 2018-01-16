package process

type Tokenizer struct{}

func NewTokenizer() (*Tokenizer, error) {
	return &Tokenizer{}, nil
}

func (tk *Tokenizer) Tokenize(file *File) error {
	file.Tokens = append(file.Tokens, file.Data)
	return nil
}
