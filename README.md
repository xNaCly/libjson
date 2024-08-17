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

| JSON size | `encoding/json` | `libjson`  |
| --------- | --------------- | ---------- |
| 64KB      | 650.457µs       | 695.546µs  |
| 128KB     | 2.689076ms      | 2.964479ms |
| 256KB     | 2.777847ms      | 3.077609ms |
| 512KB     | 5.190729ms      | 6.991226ms |

These results were generated with the following specs:

```text
OS: Arch Linux x86_64
Kernel: 6.10.4-arch2-1
Memory: 32024MiB
Go version: 1.23
```

For the first naiive implementation, these results are fairly good and not too
far behind the `encoding/go` implementation, however there are some potential
low hanging fruit for performance improvements and I will invest some time into
them, below this section I will keep a list of performance improvements and
their impact on the overall performance.

### Reproduce locally

> Make sure you have the go toolchain installed for this.

```shell
cd test/
chmod +x ./bench.sh
./bench.sh
```

Output looks something like:

```text
fetching example data
[libjson] building executable
[libjson] 64KB: 695.546µs
[libjson] 128KB: 2.964479ms
[libjson] 256KB: 3.077609ms
[libjson] 512KB: 6.991226ms
[gojson] building executable
[gojson] 64KB: 650.457µs
[gojson] 128KB: 2.689076ms
[gojson] 256KB: 2.777847ms
[gojson] 512KB: 5.190729ms
```
