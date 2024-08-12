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

- somewhat [ECMA
  404](https://ecma-international.org/wp-content/uploads/ECMA-404_2nd_edition_december_2017.pdf)
  and [rfc8259](https://www.rfc-editor.org/rfc/rfc8259) compliant
  - missing some specific edge cases :^)
- no reflection, uses a custom query language similar to JavaScript object access instead
- generics for value insertion and extraction with `libjson.Get` and `libjson.Set`
- caching of queries with `libjson.Compile`
- serialisation via `json.Marshal`
