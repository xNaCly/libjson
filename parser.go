package libjson

import (
	"errors"
	"fmt"
	"unicode/utf8"
	"unsafe"
)

type parser struct {
	l       lexer
	cur_tok token
	input   []byte
	arena   valueArena
}

func (p *parser) advance() error {
	t, err := p.l.next()
	p.cur_tok = t
	if p.cur_tok.Type == t_eof && err != nil {
		return err
	}
	return nil
}

// parses toks into a valid json representation
func (p *parser) parse(input []byte) (JSONVal, error) {
	p.input = input
	err := p.advance()
	if err != nil {
		return JSONVal{}, err
	}
	if val, err := p.expression(); err != nil {
		return JSONVal{}, err
	} else {
		if p.cur_tok.Type != t_eof {
			return JSONVal{}, fmt.Errorf("Unexpected non-whitespace character(s) (%s) after JSON data", tokennames[p.cur_tok.Type])
		}
		return val, nil
	}
}

func (p *parser) expression() (JSONVal, error) {
	if p.cur_tok.Type == t_left_curly {
		return p.object()
	} else if p.cur_tok.Type == t_left_braket {
		return p.array()
	} else {
		return p.atom()
	}
}

func (p *parser) object() (JSONVal, error) {
	if p.cur_tok.Type != t_left_curly {
		return JSONVal{}, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_left_curly])
	}
	err := p.advance()
	if err != nil {
		return JSONVal{}, err
	}

	if p.cur_tok.Type == t_right_curly {
		err := p.advance()
		if err != nil {
			return JSONVal{}, err
		}
		return p.arena.NewObjectVal(make(map[string]JSONVal, 0)), nil
	}

	m := make(map[string]JSONVal, 8)

	for p.cur_tok.Type != t_eof && p.cur_tok.Type != t_right_curly {
		if len(m) > 0 {
			if p.cur_tok.Type != t_comma {
				return JSONVal{}, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_comma])
			}
			err := p.advance()
			if err != nil {
				return JSONVal{}, err
			}
		}

		if p.cur_tok.Type != t_string {
			return JSONVal{}, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_string])
		}
		in := p.input[p.cur_tok.Start:p.cur_tok.End]
		key := *(*string)(unsafe.Pointer(&in))
		err := p.advance()
		if err != nil {
			return JSONVal{}, err
		}

		if p.cur_tok.Type != t_colon {
			return JSONVal{}, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_colon])
		}
		err = p.advance()
		if err != nil {
			return JSONVal{}, err
		}

		val, err := p.expression()
		if err != nil {
			return JSONVal{}, err
		}

		// TODO:  think about activating a uniqueness check for object keys,
		// would add an other hashing and a branch for each object key parsed.
		//
		// if _, ok := m[key]; ok {
		// 	return nil, fmt.Errorf("Key %q is already set in this object", key)
		// }

		m[key] = val
	}

	if p.cur_tok.Type != t_right_curly {
		return JSONVal{}, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_right_curly])
	}
	err = p.advance()
	if err != nil {
		return JSONVal{}, err
	}

	return p.arena.NewObjectVal(m), nil
}

func (p *parser) array() (JSONVal, error) {
	if p.cur_tok.Type != t_left_braket {
		return JSONVal{}, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_left_braket])
	}
	err := p.advance()
	if err != nil {
		return JSONVal{}, err
	}

	if p.cur_tok.Type == t_right_braket {
		return p.arena.NewArrayVal([]JSONVal{}), p.advance()
	}

	a := make([]JSONVal, 0, 8)

	for p.cur_tok.Type != t_eof && p.cur_tok.Type != t_right_braket {
		if len(a) > 0 {
			if p.cur_tok.Type != t_comma {
				return JSONVal{}, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_comma])
			}
			err := p.advance()
			if err != nil {
				return JSONVal{}, err
			}
		}
		node, err := p.expression()
		if err != nil {
			return JSONVal{}, err
		}
		a = append(a, node)
	}

	if p.cur_tok.Type != t_right_braket {
		return JSONVal{}, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_right_braket])
	}

	return p.arena.NewArrayVal(a), p.advance()
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

func (p *parser) atom() (JSONVal, error) {
	var r JSONVal
	switch p.cur_tok.Type {
	case t_string:
		in := p.input[p.cur_tok.Start:p.cur_tok.End]
		end, err := unescapeInPlace(in)
		if err != nil {
			return JSONVal{}, err
		}
		in = in[:end]
		r = p.arena.NewStringVal(*(*string)(unsafe.Pointer(&in)))
	case t_number:
		raw := p.input[p.cur_tok.Start:p.cur_tok.End]
		number, err := parseFloat(raw)
		if err != nil {
			return JSONVal{}, fmt.Errorf("Invalid floating point number %q: %w", string(raw), err)
		}
		r = NewNumberVal(number)
	case t_true:
		r = NewBoolVal(true)
	case t_false:
		r = NewBoolVal(false)
	case t_null:
		r = NewNullVal()
	default:
		return JSONVal{}, fmt.Errorf("Unexpected %q at this position, expected any of: string, number, true, false or null", tokennames[p.cur_tok.Type])
	}
	if err := p.advance(); err != nil {
		return JSONVal{}, err
	}
	return r, nil
}
