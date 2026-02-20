package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"runtime/pprof"

	"github.com/xnacly/libjson"
)

func Must[T any](t T, err error) T {
	if err != nil {
		log.Fatalln(err)
	}
	return t
}

func main() {
	noGc := flag.Bool("nogc", false, "disable the go garbage collector")
	useLibjson := flag.Bool("libjson", true, "use libjson, if false use encoding/json")
	usePprof := flag.Bool("pprof", false, "use pprof cpu tracing")
	query := flag.String("q", ".", "query the parsed json")
	silent := flag.Bool("s", false, "no stdoutput")
	escape := flag.Bool("e", false, "escapes input with Gos '%#+v'")
	flag.Parse()

	if *noGc {
		debug.SetGCPercent(-1)
	}

	args := flag.Args()

	var filePath string
	var file *os.File
	if info, err := os.Stdin.Stat(); err != nil || info.Mode()&os.ModeCharDevice != 0 { // we are in a pipe
		if len(args) == 0 {
			log.Fatalln("Wanted a file as an argument, got nothing, exiting")
		}
		filePath = args[0]
		file = Must(os.Open(filePath))
	} else {
		file = os.Stdin
		filePath = "stdin"
	}

	if *usePprof {
		f, err := os.Create(filepath.Base(filePath) + ".pprof")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *useLibjson {
		out := Must(libjson.NewReader(file))
		if !*silent {
			out := Must(libjson.Get[any](&out, *query))
			if *escape {
				fmt.Printf("%#+v\n", out)
			} else {
				fmt.Println(out)
			}
		}
	} else {
		if *query != "." {
			panic("With -libjson=false, there is no support for querying the json")
		}

		decoder := json.NewDecoder(file)
		var out any
		if err := decoder.Decode(&out); err != nil {
			panic(err)
		}

		if !*silent {
			if *escape {
				fmt.Printf("%#+v\n", out)
			} else {
				fmt.Println(out)
			}
		}
	}
}
