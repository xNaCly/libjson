package libjson

import (
	"bufio"
	"io"
	"strings"
)

func NewReader(r io.Reader) (*JSON, error) {
	p := &parser{l: lexer{r: bufio.NewReader(r)}}
	obj, err := p.parse()
	if err != nil {
		return nil, err
	}
	return &JSON{obj}, nil
}

func New(s string) (*JSON, error) {
	r := strings.NewReader(s)
	return NewReader(r)
}
