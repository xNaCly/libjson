package libjson

import (
	"bufio"
	"errors"
	"io"
)

type lexer struct {
	r *bufio.Reader
}

func (l *lexer) advanceInt(amount int) ([]byte, bool) {
	b := make([]byte, amount)
	readBytes, err := l.r.Read(b)
	return b, err == nil && readBytes == amount
}

func (l *lexer) advance() (rune, bool) {
	cc, _, err := l.r.ReadRune()
	if err != nil {
		return 0, false
	}
	return cc, true
}

func (l *lexer) lex(r io.Reader) ([]token, error) {
	l.r = bufio.NewReader(r)

	toks := make([]token, 0, l.r.Size()/2)
	for {
		cc, ok := l.advance()
		if !ok {
			break
		}

		tt := t_unkown

		switch cc {
		case ' ', '\n', '\t', '\r':
			continue
		case '{':
			tt = t_left_curly
		case '}':
			tt = t_right_curly
		case '[':
			tt = t_left_braket
		case ']':
			tt = t_right_braket
		case ',':
			tt = t_comma
		case ':':
			tt = t_colon
		case 't': // this should always be the 'true' atom and is therefore optimised here
			if r, ok := l.advanceInt(3); ok && (r[0] == 'r' && r[1] == 'u' && r[2] == 'e') {
				tt = t_true
			} else {
				return nil, errors.New("Failed to read the expected 'true' atom")
			}
		case 'f': // this should always be the 'false' atom and is therefore optimised here
			if r, ok := l.advanceInt(4); ok && (r[0] == 'a' && r[1] == 'l' && r[2] == 's' && r[3] == 'e') {
				tt = t_false
			} else {
				return nil, errors.New("Failed to read the expected 'false' atom")
			}
		case 'n': // this should always be the 'null' atom and is therefore optimised here
			if r, ok := l.advanceInt(3); ok && (r[0] == 'u' && r[1] == 'l' && r[2] == 'l') {
				tt = t_null
			} else {
				return nil, errors.New("Failed to read the expected 'null' atom")
			}

			//  TODO: numbers, strings
		}

		toks = append(toks, token{Type: tt})
	}
	return toks, nil
}
