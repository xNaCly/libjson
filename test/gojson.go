package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

func MustVal[T any](t T, err error) T {
	if err != nil {
		log.Fatalln(err)
	}
	return t
}

func Must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	args := os.Args
	if len(args) == 1 {
		log.Fatalln("Wanted a file as first argument, got nothing, exiting")
	}
	file := MustVal(os.Open(args[1]))
	m := []struct {
		Name     string
		Language string
		Id       string
		Bio      string
		Version  float64
	}{}
	s := time.Now()
	d := json.NewDecoder(file)
	Must(d.Decode(&m))
	end := time.Now().Sub(s)
	fmt.Println(end)
}
