# gojson

> WARNING: gojson is currently a work in progress :)

Fast and minimal json parser written with and for go

```go
package main

import (
    "github.com/xnacly/gojson"
)

func main() {
    input := `{ "hello": {"world": ["hi"] } }`
    jsonObj := gojson.New(jsonObj) // or gojson.NewReader(r io.Reader)

    // accessing values
    fmt.Println(gojson.Get(jsonObj, "hello.world.0")) // hi

    // updating values
    gojson.Set(jsonObj, "hello.world.0", "heyho")
    fmt.Println(gojson.Get(jsonObj, "hello.world.0")) // heyho
    gojson.Set(jsonObj, "hello.world", []string{"hi", "heyho"})
    fmt.Println(jsonObj.Get(jsonObj, "hello.world")) // []string{"hi", "heyho"}

    // compiling queries for faster access
    helloWorldQuery, _ := gojson.Compile(jsonObj, "hello.world")
    cachedQuery,  _ := helloWorldQuery()
    fmt.Println(cachedQuery)
}
```

## Features

- somewhat [ECMA
  404](https://ecma-international.org/wp-content/uploads/ECMA-404_2nd_edition_december_2017.pdf)
  and [rfc8259](https://www.rfc-editor.org/rfc/rfc8259) compliant
  - missing some specific edge cases :^)
- no reflection a custom query language similar to javascript object access instead
- generics for value insertion and extraction with `gojson.Get` and `gojson.Set`
- caching of queries with `gojson.Compile`
