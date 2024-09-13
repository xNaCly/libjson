package libjson

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectAtom(t *testing.T) {
	input := []struct {
		inp      string
		path     string
		expected any
	}{
		{"12", ".", 12},
		{`"str"`, ".", "str"},
		{"true", ".", true},
		{"false", ".", false},
		{"null", ".", nil},
		{`{"key": "value"}`, ".key", "value"},
		{`{ "hello": {"world": ["hi"] } }`, ".hello.world.0", "hi"},
	}
	for _, i := range input {
		t.Run(i.inp+i.path, func(t *testing.T) {
			obj, err := New(i.inp)
			assert.NoError(t, err)
			assert.NotNil(t, obj)
			out, err := obj.get(i.path)
			assert.NoError(t, err)
			assert.EqualValues(t, i.expected, out)
		})
	}
}

// This tests the example in the readme, always copy from here to the readme
func TestObjectReadme(t *testing.T) {
	input := `{ "hello": {"world": ["hi"] } }`
	jsonObj, _ := New(input) // or libjson.NewReader(r io.Reader)

	// accessing values
	fmt.Println(Get[string](jsonObj, ".hello.world.0")) // hi

	// updating values
	Set(jsonObj, ".hello.world.0", "heyho")
	fmt.Println(Get[string](jsonObj, ".hello.world.0")) // heyho
	Set(jsonObj, ".hello.world", []string{"hi", "heyho"})
	fmt.Println(Get[string](jsonObj, ".hello.world")) // []string{"hi", "heyho"}

	// compiling queries for faster access
	// helloWorldQuery, _ := Compile[[]any](jsonObj, ".hello.world")
	// cachedQuery, _ := helloWorldQuery()
	// fmt.Println(cachedQuery)
}

func TestStandardFail(t *testing.T) {
	input := []string{
		`{"a":"b"}/**/`,
		`{"a":"b"}/**//`,
		`{"a":"b"}//`,
		`{"a":"b"}/`,
		`{"a":"b"}#`,
	}
	for _, i := range input {
		t.Run(i, func(t *testing.T) {
			p := parser{l: lexer{r: bufio.NewReader(strings.NewReader(i))}}
			_, err := p.parse()
			assert.Error(t, err)
		})
	}
}
