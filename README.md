# libjson

> WARNING: libjson is currently a work in progress :)

Fast and minimal JSON parser written in and for go

```go
package main

import (
    "github.com/xnacly/libjson"
)

func main() {
    input := `{ "hello": {"world": ["hi"] } }`
    jsonObj := libjson.New(jsonObj) // or libjson.NewReader(r io.Reader)

    // accessing values
    fmt.Println(libjson.Get(jsonObj, "hello.world.0")) // hi

    // updating values
    libjson.Set(jsonObj, "hello.world.0", "heyho")
    fmt.Println(libjson.Get(jsonObj, "hello.world.0")) // heyho
    libjson.Set(jsonObj, "hello.world", []string{"hi", "heyho"})
    fmt.Println(jsonObj.Get(jsonObj, "hello.world")) // []string{"hi", "heyho"}

    // compiling queries for faster access
    helloWorldQuery, _ := libjson.Compile(jsonObj, "hello.world")
    cachedQuery,  _ := helloWorldQuery()
    fmt.Println(cachedQuery)
}
```

## Features

- [ECMA 404](https://ecma-international.org/wp-content/uploads/ECMA-404_2nd_edition_december_2017.pdf)
  and [rfc8259](https://www.rfc-editor.org/rfc/rfc8259) compliant
  - tests against [JSONTestSuite](https://github.com/nst/JSONTestSuite), see
    [Parsing JSON is a Minefield
    💣](https://seriot.ch/projects/parsing_json.html)in the future
  - no trailing commata, comments, `Nan` or `Infinity`
  - top level atom/skalars, like strings, numbers, true, false and null
  - uft8 support via go [rune](https://go.dev/blog/strings)
- no reflection, uses a custom query language similar to JavaScript object access instead
- generics for value insertion and extraction with `libjson.Get` and `libjson.Set`
- caching of queries with `libjson.Compile`
- serialisation via `json.Marshal`

## Benchmarks

These results were generated with the following specs:

```text
OS: Arch Linux x86_64
Kernel: 6.10.4-arch2-1
Memory: 32024MiB
Go version: 1.23
```

Below this section is a list of performance improvements and their impact on
the overall performance as well as the full results of
[test/bench.sh](test/bench.sh).

### [58d9360](https://github.com/xNaCly/libjson/commit/58d9360bae0576e761e021ee52035713206fdab1)

| JSON size | `encoding/json` | `libjson` |
| --------- | --------------- | --------- |
| 1MB       | 12.2ms          | 19.9ms    |
| 5MB       | 60.2ms          | 95.2ms    |
| 10MB      | 117.2ms         | 183.8ms   |

I had to change some things to account for issues occuring in the reading of
atoms, such as true, false and null. All of those are read by buffering the
size of chars they have and reading this buffer at once, instead of iterating
and multiple reads. This did not work correctly because i used
`(*bufio.Reader).Read`, which sometimes does not read all bytes fitting in the
buffer passed into it. Thats why these commit introduces a lot of performance
regressions.

### [e08beba](https://github.com/xNaCly/libjson/commit/e08bebada39441d9b6a20cb05251488ddce68285)

| JSON size | `encoding/json` | `libjson` |
| --------- | --------------- | --------- |
| 1MB       | 11.7ms          | 13.1ms    |
| 5MB       | 55.2ms          | 64.8ms    |

The optimisation in this commit is to no longer tokenize the whole input before
starting the parser but attaching the lexer to the parser. This allows the
parser to invoke the tokenization of the next token on demand, for instance
once the parser needs to advance. This reduces the runtime around 4ms for the
1MB input and 14ms for 5MB, resulting in a 1.33x and a 1.22x runtime reduction,
pretty good for such a simple change.

### [be686d2](https://github.com/xNaCly/libjson/commit/be686d2c85c07cdfa91295052db54001d8cd5cc8)

| JSON size | `encoding/json` | `libjson` |
| --------- | --------------- | --------- |
| 1MB       | 11.7ms          | 17.4ms    |
| 5MB       | 55.2ms          | 78.5ms    |

For the first naiive implementation, these results are fairly good and not too
far behind the `encoding/go` implementation, however there are some potential
low hanging fruit for performance improvements and I will invest some time into
them.

No specific optimisations made here, except removing the check for duplicate
object keys, because
[rfc8259](https://www.rfc-editor.org/rfc/rfc8259) says:

> When the names within an object are not
> unique, the behavior of software that receives such an object is
> unpredictable. Many implementations report the last name/value pair only.
> Other implementations report an error or fail to parse the object, and some
> implementations report all of the name/value pairs, including duplicates.

Thus I can decide wheter or not I want to error on duplicate keys, or simply
let each duplicate key overwrite the previous value in the object, however
checking if a given key is already in the map/object requires that key to be
hashed and the map to be indexed with that key, omitting this check saves us
these operations, thus making the parser faster for large objects.

### Reproduce locally

> Make sure you have the go toolchain and python3 installed for this.

```shell
cd test/
chmod +x ./bench.sh
./bench.sh
```

Output looks something like:

```text
fetching example data
building executable
Benchmark 1: ./test ./1MB.json
  Time (mean ± σ):      13.1 ms ±   0.2 ms    [User: 12.1 ms, System: 2.8 ms]
  Range (min … max):    12.7 ms …  13.8 ms    210 runs

Benchmark 2: ./test -libjson=false ./1MB.json
  Time (mean ± σ):      11.7 ms ±   0.3 ms    [User: 9.5 ms, System: 2.1 ms]
  Range (min … max):    11.1 ms …  12.7 ms    237 runs

Summary
  ./test -libjson=false ./1MB.json ran
    1.12 ± 0.03 times faster than ./test ./1MB.json
Benchmark 1: ./test ./5MB.json
  Time (mean ± σ):      64.2 ms ±   0.9 ms    [User: 79.3 ms, System: 13.1 ms]
  Range (min … max):    62.6 ms …  67.0 ms    46 runs

Benchmark 2: ./test -libjson=false ./5MB.json
  Time (mean ± σ):      55.2 ms ±   1.1 ms    [User: 51.3 ms, System: 6.3 ms]
  Range (min … max):    53.6 ms …  58.0 ms    53 runs

Summary
  ./test -libjson=false ./5MB.json ran
    1.16 ± 0.03 times faster than ./test ./5MB.json
```
