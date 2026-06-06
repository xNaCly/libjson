package libjson_test

import (
	"fmt"
	"os"

	"github.com/xnacly/libjson"
)

func ExampleNew() {
	input := `{ "hello": {"world": ["hi"] } }`

	doc, err := libjson.New([]byte(input))
	if err != nil {
		panic(err)
	}

	value, err := libjson.Get[string](&doc, ".hello.world.0")
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
	// Output: hi
}

func ExampleFromFile() {
	f, err := os.CreateTemp("", "libjson-example-*.json")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())

	_, err = f.WriteString(`{ "hello": {"world": ["hi"] } }`)
	if err != nil {
		panic(err)
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	doc, err := libjson.FromFile(f)
	if err != nil {
		panic(err)
	}
	defer doc.Close()
	defer f.Close()

	value, err := libjson.Get[string](&doc, ".hello.world.0")
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
	// Output: hi
}
