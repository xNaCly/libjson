package libjson

import (
	"errors"
	"fmt"
	"io"
)

type lexer struct {
	data []byte
	pos  int
}

func (l *lexer) advance() (byte, error) {
	if l.pos >= len(l.data) {
		return 0, io.EOF
	}
	cc := l.data[l.pos]
	l.pos++
	return cc, nil
}

func (l *lexer) next() (token, error) {
	cc, err := l.advance()
	if err != nil {
		return empty, nil
	}

	tt := t_eof

	for cc == ' ' || cc == '\n' || cc == '\t' || cc == '\r' {
		cc, err = l.advance()
		if err != nil {
			return empty, nil
		}
	}

	switch cc {
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
	case '"':
		start := l.pos
		end := start
		for {
			cc, err = l.advance()
			if cc == '"' {
				end = l.pos - 1
				break
			} else if err != nil {
				return empty, errors.New("Unterminated string detected")
			}
		}
		t := token{Type: t_string, Val: l.data[start:end]}
		return t, nil
	case 't': // this should always be the 'true' atom and is therefore optimised here
		if l.pos+3 > len(l.data) {
			return empty, errors.New("Failed to read the expected 'true' atom")
		}
		if !(l.data[l.pos] == 'r' && l.data[l.pos+1] == 'u' && l.data[l.pos+2] == 'e') {
			return empty, errors.New("Failed to read the expected 'true' atom")
		}
		l.pos += 3
		tt = t_true
	case 'f': // this should always be the 'false' atom and is therefore optimised here
		if l.pos+4 > len(l.data) {
			return empty, errors.New("Failed to read the expected 'false' atom")
		}
		if !(l.data[l.pos] == 'a' && l.data[l.pos+1] == 'l' && l.data[l.pos+2] == 's' && l.data[l.pos+3] == 'e') {
			return empty, errors.New("Failed to read the expected 'false' atom")
		}
		l.pos += 4
		tt = t_false
	case 'n': // this should always be the 'null' atom and is therefore optimised here
		if l.pos+3 > len(l.data) {
			return empty, errors.New("Failed to read the expected 'null' atom")
		}
		if !(l.data[l.pos] == 'u' && l.data[l.pos+1] == 'l' && l.data[l.pos+2] == 'l') {
			return empty, errors.New("Failed to read the expected 'null' atom")
		}
		l.pos += 3
		tt = t_null
	default:
		if cc == '-' || (cc >= '0' && cc <= '9') {
			start := l.pos - 1
			cc, err = l.advance()
			if err != nil {
				// we hit eof here
				return token{Type: t_number, Val: []byte{l.data[start]}}, nil
			}

			for {
				if (cc >= '0' && cc <= '9') || cc == '-' || cc == '+' || cc == '.' || cc == 'e' || cc == 'E' {
					cc, err = l.advance()
					if err != nil {
						break
					}
				} else {
					// the read at the start of the for loop iterates us too
					// far, if we didnt break out of the loop but exited it
					// according to its condition, thus we skip that here
					l.pos--
					break
				}
			}

			t := token{Type: t_number, Val: l.data[start:l.pos]}
			return t, nil
		} else {
			return empty, fmt.Errorf("Unexpected character %q at this position.", cc)
		}
	}

	return token{tt, nil}, nil
}

// lex is only intended for tests, use lexer.next() for production code
func (l *lexer) lex(r io.Reader) ([]token, error) {
	var err error
	l.data, err = io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	toks := make([]token, 0, len(l.data)/2)
	for {
		if tok, err := l.next(); err == nil {
			if tok.Type == t_eof {
				break
			}
			toks = append(toks, tok)
		} else {
			return nil, err
		}
	}
	if len(toks) == 0 {
		return nil, errors.New("Unexpected end of JSON input")
	}

	return toks, nil
}
