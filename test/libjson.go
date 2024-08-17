package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xnacly/libjson"
)

func MustVal[T any](t T, err error) T {
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
	file := MustVal(os.Open(args[1]))
	s := time.Now()
	MustVal(libjson.Get[any](MustVal(libjson.NewReader(file)), ""))
	end := time.Now().Sub(s)
	fmt.Println(end)
}
