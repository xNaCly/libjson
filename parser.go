package libjson

import (
	"fmt"
	"strconv"
	"unsafe"
)

type parser struct {
	l lexer
	t token
}

func (p *parser) advance() error {
	t, err := p.l.next()
	p.t = t
	if p.t.Type == t_eof && err == nil {
		return nil
	}
	return err
}

func (p *parser) expect(t t_json) error {
	if p.t.Type != t {
		return fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.t.Type], tokennames[t])
	}
	return p.advance()
}

// parses toks into a valid json representation, thus the return type can be
// either map[string]any, []any, string, nil, false, true or a number
func (p *parser) parse() (any, error) {
	err := p.advance()
	if err != nil {
		return nil, err
	}
	if val, err := p.expression(); err != nil {
		return nil, err
	} else {
		if p.t.Type != t_eof {
			return nil, fmt.Errorf("Unexpected non-whitespace character(s) (%s) after JSON data", tokennames[p.t.Type])
		}
		return val, nil
	}
}

func (p *parser) expression() (any, error) {
	if p.t.Type == t_left_curly {
		return p.object()
	} else if p.t.Type == t_left_braket {
		return p.array()
	} else {
		return p.atom()
	}
}

func (p *parser) object() (map[string]any, error) {
	err := p.expect(t_left_curly)
	if err != nil {
		return nil, err
	}

	m := make(map[string]any, 8)

	if p.t.Type == t_right_curly {
		err = p.advance()
		if err != nil {
			return nil, err
		}
		return m, nil
	}

	for p.t.Type != t_eof && p.t.Type != t_right_curly {
		if len(m) > 0 {
			err := p.expect(t_comma)
			if err != nil {
				return nil, err
			}
		}

		key := *(*string)(unsafe.Pointer(&p.t.Val))
		err := p.expect(t_string)
		if err != nil {
			return nil, err
		}

		err = p.expect(t_colon)
		if err != nil {
			return nil, err
		}

		val, err := p.expression()
		if err != nil {
			return nil, err
		}

		// TODO:  think about activating a uniqueness check for object keys,
		// would add an other hashing and a branch for each object key parsed.
		//
		// if _, ok := m[key]; ok {
		// 	return nil, fmt.Errorf("Key %q is already set in this object", key)
		// }

		m[key] = val
	}

	err = p.expect(t_right_curly)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (p *parser) array() ([]any, error) {
	err := p.expect(t_left_braket)
	if err != nil {
		return nil, err
	}

	if p.t.Type == t_right_braket {
		err = p.advance()
		return []any{}, err
	}

	a := make([]any, 0, 8)

	for p.t.Type != t_eof && p.t.Type != t_right_braket {
		if len(a) > 0 {
			err := p.expect(t_comma)
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

	return a, p.expect(t_right_braket)
}

func (p *parser) atom() (any, error) {
	var r any
	switch p.t.Type {
	case t_string:
		r = *(*string)(unsafe.Pointer(&p.t.Val))
	case t_number:
		number, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&p.t.Val)), 64)
		if err != nil {
			return empty, fmt.Errorf("Invalid floating point number %q: %w", p.t.Val, err)
		}
		r = number
	case t_true:
		r = true
	case t_false:
		r = false
	case t_null:
		r = nil
	default:
		return nil, fmt.Errorf("Unexpected %q at this position, expected any of: string, number, true, false or null", tokennames[p.t.Type])
	}
	if err := p.advance(); err != nil {
		return nil, err
	}
	return r, nil
}
