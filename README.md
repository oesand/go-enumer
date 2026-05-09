# go-enumer

`go-enumer` is a small Go library for runtime enumeration registration and lookup.
It supports two styles:

- `github.com/oesand/go-enumer` — a generic enum registry for supported scalar value types.
- `github.com/oesand/go-enumer/enumeration` — a registry for custom types that expose a stable key.

## Features

- Declare enum values for supported scalar types
- Look up enum values by name or key
- Validate values against declared enums
- Iterate over declared enum names or values
- Find values by predicate

## Installation

```bash
go get github.com/oesand/go-enumer@latest
```

## Root package usage (`enum`)

```go
package main

import (
    "fmt"
    "slices"

    "github.com/oesand/go-enumer"
)

type Status int

var (
    StatusPending = enum.DeclareNamed[Status]("Pending", 1)
    StatusActive  = enum.DeclareNamed[Status]("Active", 2)
)

func main() {
    fmt.Println(enum.IsValid[Status](1))      // true
    fmt.Println(enum.IsValid[Status](99))     // false

    value, ok := enum.FromName[Status]("Active")
    fmt.Println(value, ok)                    // 2 true

    names := slices.Collect(enum.Names[Status]())
    fmt.Println(names)

    values := slices.Collect(enum.Values[Status]())
    fmt.Println(values)
}
```

## Subpackage usage (`enumeration`)

```go
package main

import (
    "fmt"
    "slices"

    "github.com/oesand/go-enumer/enumeration"
)

type Color string

func (c Color) EnumerationKey() any {
    return string(c)
}

var (
    ColorRed   = enumeration.Declare[Color]("red")
    ColorGreen = enumeration.Declare[Color]("green")
)

func main() {
    value, ok := enumeration.FromKey[Color]("green")
    fmt.Println(value, ok)                    // green true

    found, ok := enumeration.FindValue[Color](func(c Color) bool {
        return c == "red"
    })
    fmt.Println(found, ok)                    // red true

    values := slices.Collect(enumeration.Values[Color]())
    fmt.Println(values)
}
```

## Requirements

- Go 1.23 or newer
