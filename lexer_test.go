package libjson

import (
	"strings"
	"testing"
)

func validToks(t *testing.T, expected []t_json, got []token) {
	if len(expected) != len(got) {
		t.Errorf("Lengths of 'expected' (%d) and 'got' (%d) do not match\n", len(expected), len(got))
	}

	for i, e := range expected {
		if e != got[i].Type {
			t.Errorf("Expected %#+v at %d, got %#+v", e, i, got[i])
		}
	}
}

func TestLexerWhitespace(t *testing.T) {
	json := ""
	l := lexer{}
	toks, err := l.lex(strings.NewReader(json))
	if err != nil {
		t.Error(err)
	}
	tList := []t_json{}
	validToks(t, tList, toks)
}

func TestLexerStructure(t *testing.T) {
	json := "{}[],:"
	l := lexer{}
	toks, err := l.lex(strings.NewReader(json))
	if err != nil {
		t.Error(err)
	}
	tList := []t_json{
		t_left_curly,
		t_right_curly,
		t_left_braket,
		t_right_braket,
		t_comma,
		t_colon,
	}
	validToks(t, tList, toks)
}

func TestLexerAtoms(t *testing.T) {
	json := `"string" 12345 true false null`
	l := lexer{}
	toks, err := l.lex(strings.NewReader(json))
	if err != nil {
		t.Error(err)
	}
	tList := []t_json{
		t_string,
		t_number,
		t_true,
		t_false,
		t_null,
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
	if len(toks) == 0 {
		t.Error("Not enough tokens, something went wrong")
	}

	tList := []t_json{
		t_left_braket,
		t_string,
		t_colon,
		t_string,
		t_comma,

		t_string,
		t_colon,
		t_left_braket,
		t_string,
		t_comma,
		t_number,
		t_colon,
		t_true,
		t_colon,
		t_false,
		t_colon,
		t_null,
		t_right_braket,
		t_comma,

		t_string,
		t_colon,
		t_left_curly,
		t_string,
		t_colon,
		t_string,
		t_right_curly,
		t_comma,

		t_right_curly,
	}
	if len(toks) == 0 || len(toks) != len(tList) {
		t.Error("Not enough tokens, something went wrong")
	}
}
