package libjson

import (
	"io"
	"strings"
)

func NewReader(r io.Reader) (*JSON, error) {
	toks, err := (&lexer{}).lex(r)
	obj, err := parse(toks)
	if err != nil {
		return nil, err
	}
	return &JSON{obj}, nil
}

func New(s string) (*JSON, error) {
	r := strings.NewReader(s)
	return NewReader(r)
}
