package gojson

import (
	"strings"
	"testing"
)

func TestLexerStructure(t *testing.T) {
	json := "{}[],:"
	l := lexer{}
	l.init(strings.NewReader(json))
	toks, err := l.lex()
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
	if len(toks) == 0 || len(toks) != len(tList) {
		t.Error("Not enough tokens, something went wrong")
	}
}

func TestLexerAtoms(t *testing.T) {
	json := `"string" 12345 true false null`
	l := lexer{}
	l.init(strings.NewReader(json))
	toks, err := l.lex()
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
	if len(toks) == 0 || len(toks) != len(tList) {
		t.Error("Not enough tokens, something went wrong")
	}
}

func TestLexer(t *testing.T) {
	everythingJSON := `
    {
        "key": "value",
        "arrayOfDataTypes": ["string", 1234, true, false, null],
        "subobject": { "key": "value" },
    }
    `
	l := lexer{}
	l.init(strings.NewReader(everythingJSON))
	toks, err := l.lex()
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
