package libjson

import (
	"errors"
	"fmt"
	"io"
)

type lexer struct {
	data []byte
	pos  int
	len  int
}

var numChar [256]bool

func init() {
	for c := byte('0'); c <= '9'; c++ {
		numChar[c] = true
	}
	numChar['-'] = true
	numChar['+'] = true
	numChar['.'] = true
	numChar['e'] = true
	numChar['E'] = true
}

func (l *lexer) next() (token, error) {
	for l.pos < l.len {
		cc := l.data[l.pos]
		if cc == ' ' || cc == '\n' || cc == '\t' || cc == '\r' {
			l.pos++
		} else {
			break
		}
	}

	if l.pos >= l.len {
		return empty, nil
	}

	tt := t_eof
	cc := l.data[l.pos]
	l.pos++

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
		for i := start; i < l.len; i++ {
			if c := l.data[i]; c == '"' {
				t := token{Type: t_string, Start: start, End: i}
				l.pos = i + 1
				return t, nil
			} else if c == '\\' { // OH NO ITS ESCAPING :O
				i++
				if i >= l.len {
					return empty, errors.New("Unterminated string escape")
				}

				switch l.data[i] {
				case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
					// we simply skip the escaped char, the parser has to
				case 'u':
					if i+4 > l.len {
						return empty, errors.New("Unterminated string")
					}
					i += 4
				default:
					return empty, fmt.Errorf("Invalid escape %q", l.data[i])
				}
			}
		}
		return empty, errors.New("Unterminated string")
	case 't': // this should always be the 'true' atom and is therefore optimised here
		if l.pos+3 > l.len {
			return empty, errors.New("Failed to read the expected 'true' atom")
		}
		if !(l.data[l.pos] == 'r' && l.data[l.pos+1] == 'u' && l.data[l.pos+2] == 'e') {
			return empty, errors.New("Failed to read the expected 'true' atom")
		}
		l.pos += 3
		tt = t_true
	case 'f': // this should always be the 'false' atom and is therefore optimised here
		if l.pos+4 > l.len {
			return empty, errors.New("Failed to read the expected 'false' atom")
		}
		if !(l.data[l.pos] == 'a' && l.data[l.pos+1] == 'l' && l.data[l.pos+2] == 's' && l.data[l.pos+3] == 'e') {
			return empty, errors.New("Failed to read the expected 'false' atom")
		}
		l.pos += 4
		tt = t_false
	case 'n': // this should always be the 'null' atom and is therefore optimised here
		if l.pos+3 > l.len {
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
			for l.pos < l.len && numChar[l.data[l.pos]] {
				l.pos++
			}

			return token{Type: t_number, Start: start, End: l.pos}, nil
		} else {
			return empty, fmt.Errorf("Unexpected character %q at this position.", cc)
		}
	}

	return token{Type: tt}, nil
}

// lex is only intended for tests, use lexer.next() for production code
func (l *lexer) lex(r io.Reader) ([]token, error) {
	var err error
	l.data, err = io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	l.len = len(l.data)

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
