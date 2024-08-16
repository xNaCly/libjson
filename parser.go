package libjson

import (
	"fmt"
)

type parser struct {
	toks []token
	pos  int
}

func (p *parser) atEnd() bool {
	return p.pos >= len(p.toks)
}

func (p *parser) cur() token {
	if !p.atEnd() {
		return p.toks[p.pos]
	}
	return token{Type: t_eof}
}

func (p *parser) expect(t t_json) (token, error) {
	tok := p.cur()
	if tok.Type != t {
		return token{Type: t_eof}, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[tok.Type], tokennames[t])
	}
	p.pos++
	return tok, nil
}

// parses toks into a valid json representation, thus the return type can be
// either map[string]any, []any, string, nil, false, true or a number
func (p *parser) parse() (any, error) {
	if val, err := p.expression(); err != nil {
		return nil, err
	} else {
		if !p.atEnd() {
			return nil, fmt.Errorf("Unexpected non-whitespace character(s) (%s) after JSON data of type ", tokennames[p.cur().Type])
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
		p.pos++
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

		keyStr := key.Val.(string)
		// TODO:  think about making uniqueness check for object keys configurable
		if _, ok := m[keyStr]; ok {
			return nil, fmt.Errorf("Key %q is already set in this object", keyStr)
		}
		m[keyStr] = val
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
		p.pos++
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
	cc := p.toks[p.pos]
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

	p.pos++
	return val, nil
}
