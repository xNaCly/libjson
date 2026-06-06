# libjson

Fast and minimal JSON parser written in and for Go with a JIT query language

```go
package main

import (
    "fmt"

    "github.com/xnacly/libjson"
)

func main() {
	input := `{ "hello": {"world": ["hi"] } }`
	jsonObj, _ := libjson.New([]byte(input)) // or libjson.NewReader(r io.Reader)

	// accessing values
	fmt.Println(libjson.Get[string](&jsonObj, ".hello.world.0")) // hi, nil
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
| libjson Naive         | 47,039,465  | 34,915,053  | 500,024   | 1.94x   | 2.10x fewer     |
| encoding/json Naive   | 91,059,341  | 42,744,505  | 1,050,031 | -       | -               |
| libjson Escaped       | 41,285,970  | 25,377,563  | 350,023   | 2.29x   | 3.14x fewer     |
| encoding/json Escaped | 94,677,607  | 37,544,418  | 1,100,030 | -       | -               |
| libjson Hard          | 166,323,844 | 139,915,861 | 1,400,023 | 2.43x   | 2.14x fewer     |
| encoding/json Hard    | 404,288,637 | 173,944,520 | 3,000,032 | -       | -               |

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
BenchmarkLibJson_Naive-16                     25          47039465 ns/op        34915053 B/op     500024 allocs/op
BenchmarkLibJson_Escaped-16                   30          41285970 ns/op        25377563 B/op     350023 allocs/op
BenchmarkLibJson_Hard-16                       7         166323844 ns/op        139915861 B/op   1400023 allocs/op
BenchmarkEncodingJson_Naive-16                12          91059341 ns/op        42744505 B/op    1050031 allocs/op
BenchmarkEncodingJson_Escaped-16              12          94677607 ns/op        37544418 B/op    1100030 allocs/op
BenchmarkEncodingJson_Hard-16                  3         404288637 ns/op        173944520 B/op   3000032 allocs/op
PASS
ok      github.com/xnacly/libjson       13.603s
```

### Input Source Overhead

`FromFile` avoids the `io.ReadAll` copy/allocation path and cuts constructor
allocations for file-backed inputs, at the cost of a slightly slower setup on
this machine due to the `mmap`/`munmap` syscall overhead:

| Benchmark              | ns/op       | B/op        | allocs/op |
| ---------------------- | ----------- | ----------- | --------- |
| libjson NewReader Hard | 180,711,780 | 217,042,184 | 1,400,057 |
| libjson FromFile Hard  | 211,828,341 | 134,344,372 | 1,400,025 |

### HUGE inputs

| Input size | library       | time    | faster |
| ---------- | ------------- | ------- | ------ |
| 1MB        | libjson       | 9.4ms   | 1.69x  |
|            | encoding/json | 16.0ms  |        |
| 5MB        | libjson       | 39.0ms  | 1.79x  |
|            | encoding/json | 70.0ms  |        |
| 10MB       | libjson       | 77.9ms  | 1.76x  |
|            | encoding/json | 137.1ms |        |
| 100MB      | libjson       | 719.1ms | 1.81x  |
|            | encoding/json | 1302ms  |        |

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
  Time (mean ± σ):       9.4 ms ±   0.8 ms    [User: 8.1 ms, System: 5.0 ms]
  Range (min … max):     8.8 ms …  11.3 ms    10 runs

Benchmark 2: ./test -s -libjson=false ./1MB.json
  Time (mean ± σ):      16.0 ms ±   0.9 ms    [User: 15.2 ms, System: 4.7 ms]
  Range (min … max):    15.2 ms …  18.2 ms    10 runs

Summary
  ./test -s ./1MB.json ran
    1.69 ± 0.17 times faster than ./test -s -libjson=false ./1MB.json
Benchmark 1: ./test -s ./5MB.json
  Time (mean ± σ):      39.0 ms ±   2.7 ms    [User: 44.1 ms, System: 12.1 ms]
  Range (min … max):    37.4 ms …  46.5 ms    10 runs

Benchmark 2: ./test -s -libjson=false ./5MB.json
  Time (mean ± σ):      70.0 ms ±   2.9 ms    [User: 72.2 ms, System: 9.7 ms]
  Range (min … max):    67.4 ms …  77.7 ms    10 runs

Summary
  ./test -s ./5MB.json ran
    1.79 ± 0.15 times faster than ./test -s -libjson=false ./5MB.json
Benchmark 1: ./test -s ./10MB.json
  Time (mean ± σ):      77.9 ms ±   4.5 ms    [User: 117.4 ms, System: 19.9 ms]
  Range (min … max):    72.0 ms …  86.5 ms    10 runs

Benchmark 2: ./test -s -libjson=false ./10MB.json
  Time (mean ± σ):     137.1 ms ±   3.5 ms    [User: 169.5 ms, System: 17.0 ms]
  Range (min … max):   133.8 ms … 143.8 ms    10 runs

Summary
  ./test -s ./10MB.json ran
    1.76 ± 0.11 times faster than ./test -s -libjson=false ./10MB.json
Benchmark 1: ./test -s ./100MB.json
  Time (mean ± σ):     719.1 ms ±  12.8 ms    [User: 1080.3 ms, System: 144.2 ms]
  Range (min … max):   695.7 ms … 731.9 ms    10 runs

Benchmark 2: ./test -s -libjson=false ./100MB.json
  Time (mean ± σ):      1.302 s ±  0.013 s    [User: 1.538 s, System: 0.080 s]
  Range (min … max):    1.290 s …  1.325 s    10 runs

Summary
  ./test -s ./100MB.json ran
    1.81 ± 0.04 times faster than ./test -s -libjson=false ./100MB.json
```
