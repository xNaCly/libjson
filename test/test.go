package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xnacly/libjson"
)

func main() {
	lj := flag.Bool("libjson", true, "benchmark libjson or gojson")
	flag.Parse()
	if *lj {
		args := flag.Args()
		if len(args) == 0 {
			log.Fatalln("Wanted a file as first argument, got nothing, exiting")
		}
		file, err := os.Open(args[0])
		if err != nil {
			log.Fatalln(err)
		}
		s := time.Now()
		_, err = libjson.NewReader(file)
		if err != nil {
			log.Fatalln(err)
		}
		end := time.Now().Sub(s)
		fmt.Println(end)
	} else {
		args := flag.Args()
		if len(args) == 0 {
			log.Fatalln("Wanted a file as first argument, got nothing, exiting")
		}
		file, err := os.Open(args[0])
		if err != nil {
			log.Fatalln(err)
		}
		m := []struct {
			Name     string
			Language string
			Id       string
			Bio      string
			Version  float64
		}{}
		s := time.Now()
		d := json.NewDecoder(file)
		err = d.Decode(&m)
		if err != nil {
			log.Fatalln(err)
		}
		end := time.Now().Sub(s)
		fmt.Println(end)
	}
}
