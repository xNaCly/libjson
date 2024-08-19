package libjson

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

type lexer struct {
	r   *bufio.Reader
	buf []byte
}

func (l *lexer) next() (token, error) {
	cc, err := l.r.ReadByte()
	if err != nil {
		return empty, io.EOF
	}

	tt := t_eof

	for cc == ' ' || cc == '\n' || cc == '\t' || cc == '\r' {
		cc, err = l.r.ReadByte()
		if err != nil {
			return empty, io.EOF
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
		for {
			cc, err = l.r.ReadByte()
			if cc == '"' {
				break
			} else if err != nil {
				return empty, errors.New("Unterminated string detected")
			}
			l.buf = append(l.buf, cc)
		}
		t := token{Type: t_string, Val: l.buf}
		l.buf = make([]byte, 0, 8)
		return t, nil
	case 't': // this should always be the 'true' atom and is therefore optimised here
		if b, err := l.r.ReadByte(); err != nil && b != 'r' {
			return empty, errors.New("Failed to read the expected 'true' atom")
		}
		if b, err := l.r.ReadByte(); err != nil && b != 'u' {
			return empty, errors.New("Failed to read the expected 'true' atom")
		}
		if b, err := l.r.ReadByte(); err != nil && b != 'e' {
			return empty, errors.New("Failed to read the expected 'true' atom")
		}
		tt = t_true
	case 'f': // this should always be the 'false' atom and is therefore optimised here
		if b, err := l.r.ReadByte(); err != nil && b != 'a' {
			return empty, errors.New("Failed to read the expected 'false' atom")
		}
		if b, err := l.r.ReadByte(); err != nil && b != 'l' {
			return empty, errors.New("Failed to read the expected 'false' atom")
		}
		if b, err := l.r.ReadByte(); err != nil && b != 's' {
			return empty, errors.New("Failed to read the expected 'false' atom")
		}
		if b, err := l.r.ReadByte(); err != nil && b != 'e' {
			return empty, errors.New("Failed to read the expected 'false' atom")
		}
		tt = t_false
	case 'n': // this should always be the 'null' atom and is therefore optimised here
		if b, err := l.r.ReadByte(); err != nil && b != 'u' {
			return empty, errors.New("Failed to read the expected 'null' atom")
		}
		if b, err := l.r.ReadByte(); err != nil && b != 'l' {
			return empty, errors.New("Failed to read the expected 'null' atom")
		}
		if b, err := l.r.ReadByte(); err != nil && b != 'l' {
			return empty, errors.New("Failed to read the expected 'null' atom")
		}
		tt = t_null
	default:
		if cc == '-' || (cc >= '0' && cc <= '9') {
			l.buf = append(l.buf, cc)
			cc, err = l.r.ReadByte()
			if err != nil {
				break
			}
			for {
				if (cc >= '0' && cc <= '9') || cc == '-' || cc == '+' || cc == '.' || cc == 'e' || cc == 'E' {
					l.buf = append(l.buf, cc)
					cc, err = l.r.ReadByte()
					if err != nil {
						break
					}
				} else {
					// the read at the start of the for loop iterates us too
					// far, if we didnt break out of the loop but exited it
					// according to its condition, thus we skip that here
					l.r.UnreadByte()
					break
				}
			}

			t := token{Type: t_number, Val: l.buf}
			l.buf = make([]byte, 0, 8)
			return t, nil
		} else {
			return empty, fmt.Errorf("Unexpected character %q at this position.", cc)
		}
	}

	return token{tt, nil}, nil
}

// lex is only intended for tests, use lexer.next() for production code
func (l *lexer) lex(r io.Reader) ([]token, error) {
	l.r = bufio.NewReader(r)

	toks := make([]token, 0, l.r.Size()/2)
	for {
		if tok, err := l.next(); err == nil {
			toks = append(toks, tok)
		} else {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
	}
	if len(toks) == 0 {
		return nil, errors.New("Unexpected end of JSON input")
	}

	return toks, nil
}
