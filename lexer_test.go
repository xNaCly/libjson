package libjson

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexerWhitespace(t *testing.T) {
	json := "\n\r\t      "
	l := lexer{}
	toks, err := l.lex(strings.NewReader(json))
	assert.Error(t, err)
	assert.Empty(t, toks)
}

func TestLexerStructure(t *testing.T) {
	json := "{}[],:"
	l := lexer{}
	toks, err := l.lex(strings.NewReader(json))
	if err != nil {
		t.Error(err)
	}
	tList := []token{
		{Type: t_left_curly},
		{Type: t_right_curly},
		{Type: t_left_braket},
		{Type: t_right_braket},
		{Type: t_comma},
		{Type: t_colon},
	}
	assert.EqualValues(t, tList, toks)
}

func TestLexerAtoms(t *testing.T) {
	json := `
    "string""" "🤣"
    true false null
    1 0 12.5 1e15 -1929 -0
    -1.4E+5 -129.1928e-19028
    `
	l := lexer{}
	toks, err := l.lex(strings.NewReader(json))
	assert.NoError(t, err)
	tList := []token{
		{Type: t_string, Val: []byte("string")},
		{Type: t_string, Val: []byte("")},
		{Type: t_string, Val: []byte("🤣")},
		{Type: t_true},
		{Type: t_false},
		{Type: t_null},
		{Type: t_number, Val: []byte("1")},
		{Type: t_number, Val: []byte("0")},
		{Type: t_number, Val: []byte("12.5")},
		{Type: t_number, Val: []byte("1e15")},
		{Type: t_number, Val: []byte("-1929")},
		{Type: t_number, Val: []byte("-0")},
		{Type: t_number, Val: []byte("-1.4E+5")},
		{Type: t_number, Val: []byte("-129.1928e-19028")},
	}
	assert.EqualValues(t, tList, toks)
}

func TestLexer(t *testing.T) {
	json := `
    {
        "key": "value",
        "arrayOfDataTypes": ["string", 123456789, true, false, null],
        "subobject": { "key": "value" },
    }
    `
	l := lexer{}
	toks, err := l.lex(strings.NewReader(json))
	if err != nil {
		t.Error(err)
	}

	tList := []token{
		{Type: t_left_curly},

		{Type: t_string, Val: []byte("key")},
		{Type: t_colon},
		{Type: t_string, Val: []byte("value")},
		{Type: t_comma},

		{Type: t_string, Val: []byte("arrayOfDataTypes")},
		{Type: t_colon},
		{Type: t_left_braket},
		{Type: t_string, Val: []byte("string")},
		{Type: t_comma},
		{Type: t_number, Val: []byte("123456789")},
		{Type: t_comma},
		{Type: t_true},
		{Type: t_comma},
		{Type: t_false},
		{Type: t_comma},
		{Type: t_null},
		{Type: t_right_braket},
		{Type: t_comma},

		{Type: t_string, Val: []byte("subobject")},
		{Type: t_colon},
		{Type: t_left_curly},
		{Type: t_string, Val: []byte("key")},
		{Type: t_colon},
		{Type: t_string, Val: []byte("value")},
		{Type: t_right_curly},
		{Type: t_comma},
		{Type: t_right_curly},
	}

	assert.EqualValues(t, tList, toks)

}

func TestLexerFail(t *testing.T) {
	input := []string{
		"",
		`"`,
		"'",
		`0xFF`,
		string([]byte{0x0C}),
		`{"test": 'value'}`,
		"🤣",
	}
	for _, in := range input {
		t.Run(in, func(t *testing.T) {
			l := &lexer{}
			toks, err := l.lex(strings.NewReader(in))
			assert.Error(t, err)
			assert.Empty(t, toks)
		})
	}
}
