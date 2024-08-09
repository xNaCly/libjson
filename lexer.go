package gojson

import (
	"bufio"
	"io"
)

type lexer struct {
	r *bufio.Reader
}

func (l *lexer) init(r io.Reader) {
	l.r = bufio.NewReader(r)
}

func (l *lexer) lex() ([]token, error) {
	todo()
	return nil, nil
}
