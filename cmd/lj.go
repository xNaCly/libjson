package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xnacly/libjson"
)

// TODO: support for piping data into lj

func Must[T any](t T, err error) T {
	if err != nil {
		log.Fatalln(err)
	}
	return t
}

func main() {
	args := os.Args
	if len(args) == 1 {
		log.Fatalln("Wanted a file as first argument, got nothing, exiting")
	}
	file := Must(os.Open(args[1]))
	if len(args) == 3 {
		json := Must(libjson.NewReader(file))
		query := os.Args[2]
		fmt.Printf("%+#v\n", Must(libjson.Get[any](json, query)))
	} else {
		fmt.Println(Must(libjson.Get[any](Must(libjson.NewReader(file)), "")))
	}
}
