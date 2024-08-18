package libjson

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type lexer struct {
	r *bufio.Reader
}

func (l *lexer) advance() (rune, bool) {
	cc, _, err := l.r.ReadRune()
	if err != nil {
		return 0, false
	}
	return cc, true
}

func (l *lexer) advanceByte() (byte, bool) {
	cb, err := l.r.ReadByte()
	if err != nil {
		return 0, false
	}
	return cb, true
}

func (l *lexer) next() (token, error) {
	cc, ok := l.advanceByte()
	if !ok {
		return empty, io.EOF
	}

	tt := t_eof

	for cc == ' ' || cc == '\n' || cc == '\t' || cc == '\r' {
		cc, ok = l.advanceByte()
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
		buf := strings.Builder{}
		buf.Grow(16)
		var cr rune
		for {
			cr, ok = l.advance()
			if !ok || cr == '"' {
				break
			}
			buf.WriteRune(cr)
		}
		if cr != '"' {
			return empty, errors.New("Unterminated string detected")
		}

		return token{Type: t_string, Val: buf.String()}, nil
	case 't': // this should always be the 'true' atom and is therefore optimised here
		b := make([]byte, 3)
		_, err := io.ReadFull(l.r, b[:])
		if err != nil || !(b[0] == 'r' && b[1] == 'u' && b[2] == 'e') {
			return empty, errors.New("Failed to read the expected 'true' atom")
		}
		tt = t_true
	case 'f': // this should always be the 'false' atom and is therefore optimised here
		b := make([]byte, 4)
		_, err := io.ReadFull(l.r, b[:])
		if err != nil || !(b[0] == 'a' && b[1] == 'l' && b[2] == 's' && b[3] == 'e') {
			return empty, errors.New("Failed to read the expected 'false' atom")
		}
		tt = t_false
	case 'n': // this should always be the 'null' atom and is therefore optimised here
		b := make([]byte, 3)
		_, err := io.ReadFull(l.r, b[:])
		if err != nil || !(b[0] == 'u' && b[1] == 'l' && b[2] == 'l') {
			return empty, errors.New("Failed to read the expected 'null' atom")
		}
		tt = t_null
	default:
		if cc == '-' || (cc >= '0' && cc <= '9') {
			buf := strings.Builder{}
			buf.Grow(16)
			for (cc >= '0' && cc <= '9') || cc == '-' || cc == '+' || cc == '.' || cc == 'e' || cc == 'E' {
				buf.WriteByte(cc)
				cc, ok = l.advanceByte()
				if !ok {
					break
				}
			}
			if number, err := strconv.ParseFloat(buf.String(), 64); err == nil {
				// the read at the start of the for loop iterates us too
				// far, thus we skip that here
				l.r.UnreadByte()
				return token{Type: t_number, Val: number}, nil
			} else {
				return empty, fmt.Errorf("Invalid floating point number %q: %w", buf.String(), err)
			}
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
