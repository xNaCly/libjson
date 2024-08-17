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

### Latest

| JSON size | `encoding/json` | `libjson`   |
| --------- | --------------- | ----------- |
| 64KB      | 670.796µs       | 692.776µs   |
| 128KB     | 2.700935ms      | 2.853186ms  |
| 256KB     | 2.734796ms      | 2.880656ms  |
| 512KB     | 5.492051ms      | 7.215417ms  |
| 1MB       | 10.85996ms      | 15.531744ms |
| 5MB       | 52.169431ms     | 77.415754ms |

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

> Make sure you have the go toolchain installed for this.

```shell
cd test/
chmod +x ./bench.sh
./bench.sh
```

Output looks something like:

```text
fetching example data
building executable
[libjson] 64KB: 692.776µs
[libjson] 128KB: 2.853186ms
[libjson] 256KB: 2.880656ms
[libjson] 512KB: 7.215417ms
[libjson] 1MB: 15.531744ms
[libjson] 5MB: 77.415754ms
[gojson] 64KB: 670.796µs
[gojson] 128KB: 2.700935ms
[gojson] 256KB: 2.734796ms
[gojson] 512KB: 5.492051ms
[gojson] 1MB: 10.85996ms
[gojson] 5MB: 52.169431ms
```
