package libjson

import (
	"errors"
	"fmt"
	"io"
)

type parser struct {
	l lexer
	t token
}

func (p *parser) atEnd() bool {
	return p.t.Type == t_eof
}

func (p *parser) cur() token {
	return p.t
}

func (p *parser) advance() error {
	var err error
	p.t, err = p.l.next()
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func (p *parser) expect(t t_json) (token, error) {
	tok := p.cur()
	if tok.Type != t {
		return token{Type: t_eof}, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[tok.Type], tokennames[t])
	}
	err := p.advance()
	if err != nil {
		return tok, err
	}
	return tok, nil
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
		if !p.atEnd() {
			return nil, fmt.Errorf("Unexpected non-whitespace character(s) (%s) after JSON data", tokennames[p.cur().Type])
		}
		return val, nil
	}
}

func (p *parser) expression() (any, error) {
	t := p.cur().Type
	if t == t_left_curly {
		return p.object()
	} else if t == t_left_braket {
		return p.array()
	} else {
		return p.atom()
	}
}

func (p *parser) object() (map[string]any, error) {
	_, err := p.expect(t_left_curly)
	if err != nil {
		return nil, err
	}

	m := map[string]any{}

	if p.cur().Type == t_right_curly {
		err = p.advance()
		if err != nil {
			return nil, err
		}
		return m, nil
	}

	for !p.atEnd() && p.cur().Type != t_right_curly {
		if len(m) > 0 {
			_, err := p.expect(t_comma)
			if err != nil {
				return nil, err
			}
		}

		key, err := p.expect(t_string)
		if err != nil {
			return nil, err
		}

		_, err = p.expect(t_colon)
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
		// if _, ok := m[key.Val.(string)]; ok {
		// 	return nil, fmt.Errorf("Key %q is already set in this object", keyStr)
		// }

		m[key.Val.(string)] = val
	}

	_, err = p.expect(t_right_curly)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (p *parser) array() ([]any, error) {
	_, err := p.expect(t_left_braket)
	if err != nil {
		return nil, err
	}
	if p.cur().Type == t_right_braket {
		err = p.advance()
		if err != nil {
			return nil, err
		}
		return []any{}, nil
	}

	a := make([]any, 0, 16)

	for !p.atEnd() && p.cur().Type != t_right_braket {
		if len(a) > 0 {
			_, err := p.expect(t_comma)
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

	_, err = p.expect(t_right_braket)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (p *parser) atom() (any, error) {
	cc := p.cur()
	var val any
	switch cc.Type {
	case t_string, t_number:
		val = cc.Val
	case t_true:
		val = true
	case t_false:
		val = false
	case t_null:
		val = nil
	default:
		return nil, fmt.Errorf("Unexpected %q at this position, expected any of: string, number, true, false or null", tokennames[cc.Type])
	}

	err := p.advance()
	if err != nil {
		return nil, err
	}

	return val, nil
}
