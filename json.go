package gojson

import (
	"io"
	"strings"
)

func todo() {
	panic("todo")
}

func NewReader(r io.Reader) (*JSON, error) {
	l := lexer{}
	l.init(r)
	toks, err := l.lex()
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
