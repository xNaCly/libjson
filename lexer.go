package libjson

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unsafe"
)

type lexer struct {
	r   *bufio.Reader
	buf []byte
}

func (l *lexer) advance() (byte, bool) {
	cc, err := l.r.ReadByte()
	if err != nil {
		return 0, false
	}
	return cc, true
}

func (l *lexer) next() (token, error) {
	cc, ok := l.advance()
	if !ok {
		return empty, io.EOF
	}

	tt := t_eof

	for cc == ' ' || cc == '\n' || cc == '\t' || cc == '\r' {
		cc, ok = l.advance()
		if !ok {
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
			cc, ok = l.advance()
			if cc == '"' {
				break
			} else if !ok {
				return empty, errors.New("Unterminated string detected")
			}
			l.buf = append(l.buf, cc)
		}
		t := token{Type: t_string, Val: *(*string)(unsafe.Pointer(&l.buf))}
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
			broke := false
			for (cc >= '0' && cc <= '9') || cc == '-' || cc == '+' || cc == '.' || cc == 'e' || cc == 'E' {
				l.buf = append(l.buf, cc)
				cc, ok = l.advance()
				if !ok {
					broke = true
					break
				}
			}

			if !broke {
				// the read at the start of the for loop iterates us too
				// far, if we didnt break out of the loop but exited it
				// according to its condition, thus we skip that here
				l.r.UnreadByte()
			}
			t := token{Type: t_number, Val: *(*string)(unsafe.Pointer(&l.buf))}
			l.buf = make([]byte, 0, 8)
			return t, nil
		} else {
			return empty, fmt.Errorf("Unexpected character %q at this position.", cc)
		}
	}

	return token{tt, ""}, nil
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
