# libjson

Fast and minimal JSON parser written in and for Go with a JIT query language

```go
package main

import (
    "github.com/xnacly/libjson"
)

func main() {
	input := `{ "hello": {"world": ["hi"] } }`
	jsonObj, _ := libjson.New([]byte(input)) // or libjson.NewReader(r io.Reader)

	// accessing values
	fmt.Println(libjson.Get[string](jsonObj, ".hello.world.0")) // hi, nil
}
```

## Features

- Parser consumes and mutates the input to make most operations zero copy and zero alloc
- Full materialisation, no type access helpers or other weird overhead 
- [ECMA 404](https://ecma-international.org/publications-and-standards/standards/ecma-404/)
  and [rfc8259](https://www.rfc-editor.org/rfc/rfc8259) compliant
  - tests against [JSONTestSuite](https://github.com/nst/JSONTestSuite), see
    [Parsing JSON is a Minefield
    💣](https://seriot.ch/projects/parsing_json.html)in the future
  - no trailing commata, comments, `Nan` or `Infinity`
  - top level atom/skalars, like strings, numbers, true, false and null
  - uft8 support via go [rune](https://go.dev/blog/strings)
- no reflection, uses a custom query language similar to JavaScript object access instead, or simply use the go values as is
- generics for value insertion and extraction with `libjson.Get` and `libjson.Set`
- caching of queries with `libjson.Compile`, just in time caching of queries
- serialisation via `json.Marshal`

## Why is it faster than encoding/json?

- zero-copy strings
- mutate input for string escaping instead of allocating
- no allocations for strings, views into the original input
- no reflection
- no copies for map keys
- very simple lexer and parser

## Benchmarks

### Go internal 

| Benchmark             | ns/op       | B/op        | allocs/op | speedup | alloc reduction |
| --------------------- | ----------- | ----------- | --------- | ------- | --------------- |
| libjson Naive         | 46,294,122  | 34,907,845  | 500,023   | 1.85x   | 2.10x fewer     |
| encoding/json Naive   | 85,502,921  | 42,744,522  | 1,050,031 | -       | -               |
| libjson Escaped       | 38,199,760  | 25,394,245  | 350,023   | 2.32x   | 3.14x fewer     |
| encoding/json Escaped | 88,478,499  | 37,544,406  | 1,100,030 | -       | -               |
| libjson Hard          | 154,178,081 | 139,915,859 | 1,400,023 | 2.52x   | 2.14x fewer     |
| encoding/json Hard    | 388,198,395 | 173,944,514 | 3,000,032 | -       | -               |

Run via

```shell
go test -bench=. -benchmem
```

Results in:

```text
goos: linux
goarch: amd64
pkg: github.com/xnacly/libjson
cpu: AMD Ryzen 7 3700X 8-Core Processor
BenchmarkLibJson_Naive-16                     26          46294122 ns/op        34907845 B/op     500023 allocs/op
BenchmarkLibJson_Escaped-16                   28          38199760 ns/op        25394245 B/op     350023 allocs/op
BenchmarkLibJson_Hard-16                       7         154178081 ns/op        139915859 B/op   1400023 allocs/op
BenchmarkEncodingJson_Naive-16                13          85502921 ns/op        42744522 B/op    1050031 allocs/op
BenchmarkEncodingJson_Escaped-16              12          88478499 ns/op        37544406 B/op    1100030 allocs/op
BenchmarkEncodingJson_Hard-16                  3         388198395 ns/op        173944514 B/op   3000032 allocs/op
PASS
ok      github.com/xnacly/libjson       8.510s
```

### HUGE inputs

| Input size | library       | time    | faster |
| ----- | ------------- | ------- | ------ |
| 1MB   | libjson       | 8.7ms   | 1.73x  |
|       | encoding/json | 15.0ms  |        |
| 5MB   | libjson       | 33.2ms  | 1.99x  |
|       | encoding/json | 66.3ms |        |
| 10MB  | libjson       | 64.4ms  | 2.04x  |
|       | encoding/json | 131.6ms |        |
| 100MB  | libjson       | 618.2ms  | 2.06x  |
|       | encoding/json | 1273ms |        |

> Make sure you have the go toolchain and python3 installed for this.

```shell
cd benchmarks/
chmod +x ./bench.sh
./bench.sh
```

Output looks something like:

```text
generating example data
building executable
Benchmark 1: ./test -s ./1MB.json
  Time (mean ± σ):       8.6 ms ±   0.2 ms    [User: 10.1 ms, System: 2.8 ms]
  Range (min … max):     8.3 ms …   8.8 ms    10 runs

Benchmark 2: ./test -s -libjson=false ./1MB.json
  Time (mean ± σ):      15.1 ms ±   0.3 ms    [User: 15.6 ms, System: 3.2 ms]
  Range (min … max):    14.7 ms …  15.6 ms    10 runs

Summary
  ./test -s ./1MB.json ran
    1.76 ± 0.05 times faster than ./test -s -libjson=false ./1MB.json
Benchmark 1: ./test -s ./5MB.json
  Time (mean ± σ):      33.6 ms ±   0.8 ms    [User: 40.4 ms, System: 10.1 ms]
  Range (min … max):    32.5 ms …  34.9 ms    10 runs

Benchmark 2: ./test -s -libjson=false ./5MB.json
  Time (mean ± σ):      66.2 ms ±   0.7 ms    [User: 66.5 ms, System: 9.4 ms]
  Range (min … max):    65.3 ms …  67.7 ms    10 runs

Summary
  ./test -s ./5MB.json ran
    1.97 ± 0.05 times faster than ./test -s -libjson=false ./5MB.json
Benchmark 1: ./test -s ./10MB.json
  Time (mean ± σ):      64.3 ms ±   1.4 ms    [User: 83.6 ms, System: 12.4 ms]
  Range (min … max):    62.9 ms …  67.5 ms    10 runs

Benchmark 2: ./test -s -libjson=false ./10MB.json
  Time (mean ± σ):     132.4 ms ±   1.4 ms    [User: 169.4 ms, System: 11.7 ms]
  Range (min … max):   130.7 ms … 135.3 ms    10 runs

Summary
  ./test -s ./10MB.json ran
    2.06 ± 0.05 times faster than ./test -s -libjson=false ./10MB.json
Benchmark 1: ./test -s ./100MB.json
  Time (mean ± σ):     613.2 ms ±   2.9 ms    [User: 803.8 ms, System: 65.6 ms]
  Range (min … max):   609.0 ms … 618.7 ms    10 runs

Benchmark 2: ./test -s -libjson=false ./100MB.json
  Time (mean ± σ):      1.276 s ±  0.012 s    [User: 1.522 s, System: 0.072 s]
  Range (min … max):    1.262 s …  1.299 s    10 runs

Summary
  ./test -s ./100MB.json ran
    2.08 ± 0.02 times faster than ./test -s -libjson=false ./100MB.json
```
