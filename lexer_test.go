package libjson

import (
	"strings"
	"testing"
)

func validToks(t *testing.T, expected []token, got []token) {
	if len(expected) != len(got) {
		t.Errorf("Lengths of 'expected' (%d) and 'got' (%d) do not match\n", len(expected), len(got))
	}

	for i, e := range expected {
		if e.Type != got[i].Type {
			t.Errorf("Expected %#+v at %d, got %#+v", e, i, got[i])
		}
		if e.Type == t_number && e.Val.(float64) != got[i].Val.(float64) {
			t.Errorf("Expected %#+v at %d, got %#+v", e, i, got[i])
		}
		if e.Type == t_string && e.Val.(string) != got[i].Val.(string) {
			t.Errorf("Expected %#+v at %d, got %#+v", e, i, got[i])
		}
	}
}

func TestLexerWhitespace(t *testing.T) {
	json := "\n\r\t      "
	l := lexer{}
	toks, err := l.lex(strings.NewReader(json))
	if err != nil {
		t.Error(err)
	}
	tList := []token{}
	validToks(t, tList, toks)
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
	validToks(t, tList, toks)
}

func TestLexerAtoms(t *testing.T) {
	json := `
    "string""" ""
    true false null
    1 0 12.5 1e15 -1929 -0
    -1.4E+5 -129.1928e-19028
    `
	l := lexer{}
	toks, err := l.lex(strings.NewReader(json))
	if err != nil {
		t.Error(err)
	}
	tList := []token{
		{Type: t_string, Val: "string"},
		{Type: t_string, Val: ""},
		{Type: t_string, Val: ""},
		{Type: t_true},
		{Type: t_false},
		{Type: t_null},
		{Type: t_number, Val: 1.0},
		{Type: t_number, Val: 0.0},
		{Type: t_number, Val: 12.5},
		{Type: t_number, Val: 1e15},
		{Type: t_number, Val: -1929.0},
		{Type: t_number, Val: -0.0},
		{Type: t_number, Val: -1.4e+5},
		{Type: t_number, Val: -129.1928e-19028},
	}
	validToks(t, tList, toks)
}

func TestLexer(t *testing.T) {
	json := `
    {
        "key": "value",
        "arrayOfDataTypes": ["string", 1234, true, false, null],
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

		{Type: t_string, Val: "key"},
		{Type: t_colon},
		{Type: t_string, Val: "value"},
		{Type: t_comma},

		{Type: t_string, Val: "arrayOfDataTypes"},
		{Type: t_colon},
		{Type: t_left_braket},
		{Type: t_string, Val: "string"},
		{Type: t_comma},
		{Type: t_number, Val: 1234.0},
		{Type: t_comma},
		{Type: t_true},
		{Type: t_comma},
		{Type: t_false},
		{Type: t_comma},
		{Type: t_null},
		{Type: t_right_braket},
		{Type: t_comma},

		{Type: t_string, Val: "subobject"},
		{Type: t_colon},
		{Type: t_left_curly},
		{Type: t_string, Val: "key"},
		{Type: t_colon},
		{Type: t_string, Val: "value"},
		{Type: t_right_curly},
		{Type: t_comma},
		{Type: t_right_curly},
	}

	validToks(t, tList, toks)
}
