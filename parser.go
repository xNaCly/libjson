package libjson

import (
	"errors"
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
	return p.toks[p.pos]
}

func (p *parser) next() (token, bool) {
	if p.pos+1 >= len(p.toks) {
		return token{}, false
	}
	return p.toks[p.pos+1], true
}

func (p *parser) expect(t t_json) error {
	if p.cur().Type != t {
		return fmt.Errorf("Wanted %q, got %q", tokennames[t], tokennames[p.toks[p.pos].Type])
	}
	p.pos++
	return nil
}

// parses toks into a valid json representation, thus the return type can be
// either map[string]any, []any, string, nil, false, true or a number
func (p *parser) parse() (any, error) {
	r := []any{}
	for p.pos < len(p.toks) {
		if val, err := p.expression(); err != nil {
			return nil, err
		} else {
			r = append(r, val)
		}
	}
	if len(r) == 1 {
		return r[0], nil
	}
	return r, nil
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
	return nil, errors.ErrUnsupported
}

func (p *parser) array() ([]any, error) {
	err := p.expect(t_left_braket)
	if err != nil {
		return nil, err
	}
	a := make([]any, 0)
	if p.cur().Type == t_right_braket {
		p.pos++
		return a, nil
	}

	for !p.atEnd() && p.cur().Type != t_right_braket {
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

	err = p.expect(t_right_braket)
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
		return nil, fmt.Errorf("Unknown atom %q at this position", tokennames[cc.Type])
	}

	p.pos++
	return val, nil
}
