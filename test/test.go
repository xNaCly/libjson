package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	// "runtime/pprof"

	"github.com/xnacly/libjson"
)

func main() {
	// f, err := os.Create("cpu.pprof")
	// if err != nil {
	// 	panic(err)
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()
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
		_, err = libjson.NewReader(file)
		if err != nil {
			log.Fatalln(err)
		}
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
		d := json.NewDecoder(file)
		err = d.Decode(&m)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
