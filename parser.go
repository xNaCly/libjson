package libjson

import (
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"
	"unsafe"
)

type parser struct {
	l       lexer
	cur_tok token
	input   []byte
}

func (p *parser) advance() error {
	t, err := p.l.next()
	p.cur_tok = t
	if p.cur_tok.Type == t_eof && err != nil {
		return err
	}
	return nil
}

// parses toks into a valid json representation, thus the return type can be
// either map[string]any, []any, string, nil, false, true or a number
func (p *parser) parse(input []byte) (any, error) {
	p.input = input
	err := p.advance()
	if err != nil {
		return nil, err
	}
	if val, err := p.expression(); err != nil {
		return nil, err
	} else {
		if p.cur_tok.Type != t_eof {
			return nil, fmt.Errorf("Unexpected non-whitespace character(s) (%s) after JSON data", tokennames[p.cur_tok.Type])
		}
		return val, nil
	}
}

func (p *parser) expression() (any, error) {
	if p.cur_tok.Type == t_left_curly {
		return p.object()
	} else if p.cur_tok.Type == t_left_braket {
		return p.array()
	} else {
		return p.atom()
	}
}

type field struct {
	Key   string
	Value any
}

type obj struct {
	Fields []field
}

var emptyObj = obj{}

func (p *parser) object() (obj, error) {
	if p.cur_tok.Type != t_left_curly {
		return emptyObj, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_left_curly])
	}
	err := p.advance()
	if err != nil {
		return emptyObj, err
	}

	if p.cur_tok.Type == t_right_curly {
		return emptyObj, p.advance()
	}

	m := obj{
		Fields: make([]field, 0, 8),
	}

	for p.cur_tok.Type != t_eof && p.cur_tok.Type != t_right_curly {
		if len(m.Fields) > 0 {
			if p.cur_tok.Type != t_comma {
				return emptyObj, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_comma])
			}
			err := p.advance()
			if err != nil {
				return emptyObj, err
			}
		}

		if p.cur_tok.Type != t_string {
			return emptyObj, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_string])
		}
		in := p.input[p.cur_tok.Start:p.cur_tok.End]
		key := *(*string)(unsafe.Pointer(&in))
		err := p.advance()
		if err != nil {
			return emptyObj, err
		}

		if p.cur_tok.Type != t_colon {
			return emptyObj, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_colon])
		}
		err = p.advance()
		if err != nil {
			return emptyObj, err
		}

		val, err := p.expression()
		if err != nil {
			return emptyObj, err
		}

		m.Fields = append(m.Fields, field{key, val})
	}

	if p.cur_tok.Type != t_right_curly {
		return emptyObj, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_right_curly])
	}
	return m, p.advance()
}

func (p *parser) array() ([]any, error) {
	if p.cur_tok.Type != t_left_braket {
		return nil, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_left_braket])
	}
	err := p.advance()
	if err != nil {
		return nil, err
	}

	if p.cur_tok.Type == t_right_braket {
		return []any{}, p.advance()
	}

	a := make([]any, 0, 8)

	for p.cur_tok.Type != t_eof && p.cur_tok.Type != t_right_braket {
		if len(a) > 0 {
			if p.cur_tok.Type != t_comma {
				return nil, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_comma])
			}
			err := p.advance()
			if err != nil {
				return nil, err
			}
		}
		node, err := p.expression()
		if err != nil {
			return nil, err
		}
		a = append(a, node)
	}

	if p.cur_tok.Type != t_right_braket {
		return nil, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_right_braket])
	}

	return a, p.advance()
}

var badEscapeErr = errors.New("bad escape")

// unescapes JSON escapes in a buffer into their non-JSON representation
//
// Returns the end of the in place escaped buffer so the caller can resize to
// the new, smaller buffer size
//
// The implementation may look weird, but is optimised to have the least
// possible branches
func unescapeInPlace(in []byte) (int, error) {
	curEnd := 0
	for i := 0; i < len(in); i++ {
		b := in[i]
		if b != '\\' {
			in[curEnd] = b
			curEnd++
			continue
		}

		// check if there’s at least 1 more byte for the escape
		if i+1 >= len(in) {
			return 0, badEscapeErr
		}
		i++ // skip \
		b = in[i]

		switch b {
		case '"', '\\', '/':
			in[curEnd] = b
			curEnd++
		case 'b':
			in[curEnd] = '\b'
			curEnd++
		case 'f':
			in[curEnd] = '\f'
			curEnd++
		case 'n':
			in[curEnd] = '\n'
			curEnd++
		case 'r':
			in[curEnd] = '\r'
			curEnd++
		case 't':
			in[curEnd] = '\t'
			curEnd++
		case 'u': // \uXXXX

			// From ECMA-404:
			//
			// However, whether a processor of JSON texts interprets such a surrogate pair
			// as a single code point or as an explicit surrogate pair is a semantic
			// decision that is determined by the specific processor.
			//
			// meaning we dont merge unicode points, firstly because fuck
			// utf16, and secondly because its simpler to just keep two unicode
			// points separate compared to increasing the complexity of this
			// decoding

			if i+4 >= len(in) {
				return 0, badEscapeErr
			}

			r, err := hex4(in[i+1 : i+5])
			if err != nil {
				return 0, err
			}
			n := utf8.EncodeRune(in[curEnd:], r)
			curEnd += n
			i += 4
		} // we dont need a default case since we check all possible escapes in the lexer
	}

	return curEnd, nil
}

func (p *parser) atom() (any, error) {
	var r any
	switch p.cur_tok.Type {
	case t_string:
		in := p.input[p.cur_tok.Start:p.cur_tok.End]
		end, err := unescapeInPlace(in)
		if err != nil {
			return nil, err
		}
		in = in[:end]
		r = *(*string)(unsafe.Pointer(&in))
	case t_number:
		raw := p.input[p.cur_tok.Start:p.cur_tok.End]
		number, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&raw)), 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid floating point number %q: %w", string(raw), err)
		}
		r = number
	case t_true:
		r = true
	case t_false:
		r = false
	case t_null:
		r = nil
	default:
		return nil, fmt.Errorf("Unexpected %q at this position, expected any of: string, number, true, false or null", tokennames[p.cur_tok.Type])
	}
	if err := p.advance(); err != nil {
		return nil, err
	}
	return r, nil
}
