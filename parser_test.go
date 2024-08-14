package libjson

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserAtoms(t *testing.T) {
	input := []string{
		"true",
		"false",
		"null",
		"12345",
		`"isastring"`,
	}
	wanted := []any{
		true,
		false,
		nil,
		12345.0,
		"isastring",
	}
	for i, in := range input {
		l := &lexer{}
		toks, err := l.lex(strings.NewReader(in))
		assert.NoError(t, err)
		p := &parser{toks, 0}
		out, err := p.parse()
		assert.NoError(t, err)
		assert.EqualValues(t, wanted[i], out)
	}
}

func TestParserArray(t *testing.T) {
	input := []string{
		"[]",
		"[1]",
		"[1, 2,3]",
		`["ayo", true, false, null, 12e12]`,
	}
	wanted := [][]any{
		{},
		{1.0},
		{1.0, 2.0, 3.0},
		{"ayo", true, false, nil, 12e12},
	}
	for i, in := range input {
		l := &lexer{}
		toks, err := l.lex(strings.NewReader(in))
		assert.NoError(t, err)
		p := &parser{toks, 0}
		out, err := p.parse()
		assert.NoError(t, err)
		assert.EqualValues(t, wanted[i], out)
	}
}
